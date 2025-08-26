package slodb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"openPAQ/internal/algorithms"
	"openPAQ/internal/normalization"
	"openPAQ/internal/types"

	_ "github.com/mattn/go-sqlite3"
)

type SIAddressDB struct {
	db         *sql.DB
	config     algorithms.MatchSeverityConfig
	normalizer *normalization.Normalizer
}

func NewSIAddressDB(dbPath string, config algorithms.MatchSeverityConfig, normalizer *normalization.Normalizer) (*SIAddressDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &SIAddressDB{db: db, config: config, normalizer: normalizer}, nil
}

func (s *SIAddressDB) Handle(ctx context.Context, input types.Input) <-chan types.PairMatching {
	c := make(chan types.PairMatching, 1)
	go func() {
		defer close(c)
		ni := input.Normalize()

		// Run three checks similar to Nominatim implementation
		cityStreetC := s.CityStreetCheck(ctx, ni)
		var cityStreetRes types.PairMatching = <-cityStreetC

		postalCityC := s.PostalCodeCityCheck(ctx, ni, cityStreetRes.StreetCityMatches)
		postalStreetC := s.PostalCodeStreetCheck(ctx, ni, cityStreetRes.StreetCityMatches)

		var result types.PairMatching
		result.StreetCityMatch = cityStreetRes.StreetCityMatch
		result.StreetCityMatches = cityStreetRes.StreetCityMatches

		for i := 0; i < 2; i++ {
			select {
			case p := <-postalCityC:
				postalCityC = nil
				result.CityPostalCodeMatch = p.CityPostalCodeMatch
				result.CityPostalCodeMatches = p.CityPostalCodeMatches
			case sres := <-postalStreetC:
				postalStreetC = nil
				result.PostalCodeStreetMatch = sres.PostalCodeStreetMatch
				result.PostalCodeStreetMatches = sres.PostalCodeStreetMatches
			case <-ctx.Done():
				c <- result
				return
			}
		}
		c <- result
	}()
	return c
}

func (s *SIAddressDB) CityStreetCheck(ctx context.Context, input types.NormalizeInput) chan types.PairMatching {
	reChan := make(chan types.PairMatching, 1)
	go func() {
		defer close(reChan)
		var matches []types.CityStreetPostalCode

		for _, street := range input.Streets {
			// Query by city and street candidates (case-insensitive)
			rows, err := s.db.QueryContext(ctx, `
				SELECT LOWER(naselje_naziv), LOWER(ulica_naziv), postni_okolis_sifra
				FROM slovenian_addresses
				WHERE LOWER(naselje_naziv) = ? AND LOWER(ulica_naziv) = ?
				GROUP BY naselje_naziv, ulica_naziv, postni_okolis_sifra
			`, strings.ToLower(input.City), strings.ToLower(street))
			if err == nil {
				for rows.Next() {
					var city, dbStreet, pc string
					if err := rows.Scan(&city, &dbStreet, &pc); err == nil {
						// Fuzzy match street and city
						validateStreet := s.config
						validateStreet.AllowPartialMatch = true
						validateStreet.AllowPartialCompareListMatch = true
						streetMatches, e1 := algorithms.GetMatches(street, []string{dbStreet}, validateStreet)
						validateCity := s.config
						validateCity.AllowPartialMatch = true
						validateCity.AllowPartialCompareListMatch = true
						cityMatches, e2 := algorithms.GetMatches(input.City, []string{city}, validateCity)
						if e1 == nil && e2 == nil {
							for _, sm := range streetMatches {
								for _, cm := range cityMatches {
									if sm.Similarity > 0 && cm.Similarity > 0 {
										matches = append(matches, types.CityStreetPostalCode{
											City:                  cm.Value,
											Street:                sm.Value,
											PostalCode:            pc,
											CountryCode:           "si",
											StreetSimilarity:      sm.Similarity,
											WasPartialStreetMatch: sm.WasPartial,
											CitySimilarity:        cm.Similarity,
											WasPartialCityMatch:   cm.WasPartial,
										})
									}
								}
							}
						}
					}
					rows.Close()
				}
			}

			// If no exact equality hits, try LIKE based narrowing to reduce set
			if len(matches) == 0 {
				rows, err := s.db.QueryContext(ctx, `
					SELECT LOWER(naselje_naziv), LOWER(ulica_naziv), postni_okolis_sifra
					FROM slovenian_addresses
					WHERE LOWER(naselje_naziv) LIKE ? AND LOWER(ulica_naziv) LIKE ?
					GROUP BY naselje_naziv, ulica_naziv, postni_okolis_sifra
				`, like(input.City), like(street))
				if err == nil {
					for rows.Next() {
						var city, dbStreet, pc string
						if err := rows.Scan(&city, &dbStreet, &pc); err == nil {
							validateStreet := s.config
							validateStreet.AllowPartialMatch = true
							validateStreet.AllowPartialCompareListMatch = true
							streetMatches, e1 := algorithms.GetMatches(street, []string{dbStreet}, validateStreet)
							validateCity := s.config
							validateCity.AllowPartialMatch = true
							validateCity.AllowPartialCompareListMatch = true
							cityMatches, e2 := algorithms.GetMatches(input.City, []string{city}, validateCity)
							if e1 == nil && e2 == nil {
								for _, sm := range streetMatches {
									for _, cm := range cityMatches {
										if sm.Similarity > 0 && cm.Similarity > 0 {
											matches = append(matches, types.CityStreetPostalCode{
												City:                  cm.Value,
												Street:                sm.Value,
												PostalCode:            pc,
												CountryCode:           "si",
												StreetSimilarity:      sm.Similarity,
												WasPartialStreetMatch: sm.WasPartial,
												CitySimilarity:        cm.Similarity,
												WasPartialCityMatch:   cm.WasPartial,
											})
										}
									}
								}
							}
						}
						rows.Close()
					}
				}
			}
		}

		unique := types.RemoveDuplicate(matches)
		res := types.PairMatching{}
		if len(unique) > 0 {
			res.StreetCityMatch = true
			res.StreetCityMatches = append(res.StreetCityMatches, unique...)
		}
		reChan <- res
	}()
	return reChan
}

func (s *SIAddressDB) PostalCodeStreetCheck(ctx context.Context, input types.NormalizeInput, cityStreet []types.CityStreetPostalCode) chan types.PairMatching {
	reChan := make(chan types.PairMatching, 1)
	go func() {
		defer close(reChan)
		var result types.PairMatching

		for _, elem := range cityStreet {
			if (strings.Contains(input.PostalCode, elem.PostalCode) && elem.PostalCode != "") || (strings.Contains(elem.PostalCode, input.PostalCode) && input.PostalCode != "") {
				result.PostalCodeStreetMatch = true
				result.PostalCodeStreetMatches = append(result.PostalCodeStreetMatches, types.PostalCodeStreet{
					PostalCode:            elem.PostalCode,
					Street:                elem.Street,
					CountryCode:           "si",
					StreetSimilarity:      elem.StreetSimilarity,
					WasPartialStreetMatch: elem.WasPartialStreetMatch,
				})
			}
		}
		if result.PostalCodeStreetMatch {
			reChan <- result
			return
		}

		for _, street := range input.Streets {
			rows, err := s.db.QueryContext(ctx, `
				SELECT LOWER(ulica_naziv), postni_okolis_sifra
				FROM slovenian_addresses
				WHERE postni_okolis_sifra = ? AND LOWER(ulica_naziv) LIKE ?
				GROUP BY ulica_naziv, postni_okolis_sifra
			`, input.PostalCode, like(street))
			if err == nil {
				for rows.Next() {
					var dbStreet, pc string
					if err := rows.Scan(&dbStreet, &pc); err == nil {
						validateStreet := s.config
						validateStreet.AllowPartialMatch = true
						validateStreet.AllowPartialCompareListMatch = true
						streetMatches, e1 := algorithms.GetMatches(street, []string{dbStreet}, validateStreet)
						if e1 == nil {
							for _, sm := range streetMatches {
								result.PostalCodeStreetMatches = append(result.PostalCodeStreetMatches, types.PostalCodeStreet{
									PostalCode:            pc,
									Street:                sm.Value,
									CountryCode:           "si",
									StreetSimilarity:      sm.Similarity,
									WasPartialStreetMatch: sm.WasPartial,
								})
							}
						}
					}
				}
				rows.Close()
			}
		}

		if len(result.PostalCodeStreetMatches) > 0 {
			result.PostalCodeStreetMatch = true
		}
		reChan <- result
	}()
	return reChan
}

func (s *SIAddressDB) PostalCodeCityCheck(ctx context.Context, input types.NormalizeInput, cityStreet []types.CityStreetPostalCode) chan types.PairMatching {
	reChan := make(chan types.PairMatching, 1)
	go func() {
		defer close(reChan)
		var result types.PairMatching

		for _, elem := range cityStreet {
			if (strings.Contains(input.PostalCode, elem.PostalCode) && elem.PostalCode != "") || (strings.Contains(elem.PostalCode, input.PostalCode) && input.PostalCode != "") {
				result.CityPostalCodeMatch = true
				result.CityPostalCodeMatches = append(result.CityPostalCodeMatches, types.CityPostalCode{
					PostalCode:          elem.PostalCode,
					City:                elem.City,
					CountryCode:         "si",
					CitySimilarity:      elem.CitySimilarity,
					WasPartialCityMatch: elem.WasPartialCityMatch,
				})
			}
		}
		if result.CityPostalCodeMatch {
			reChan <- result
			return
		}

		rows, err := s.db.QueryContext(ctx, `
			SELECT LOWER(naselje_naziv), postni_okolis_sifra
			FROM slovenian_addresses
			WHERE postni_okolis_sifra = ?
			GROUP BY naselje_naziv, postni_okolis_sifra
		`, input.PostalCode)
		if err == nil {
			for rows.Next() {
				var city, pc string
				if err := rows.Scan(&city, &pc); err == nil {
					validateCity := s.config
					validateCity.AllowPartialMatch = true
					validateCity.AllowPartialCompareListMatch = true
					cityMatches, e1 := algorithms.GetMatches(input.City, []string{city}, validateCity)
					if e1 == nil {
						for _, cm := range cityMatches {
							result.CityPostalCodeMatches = append(result.CityPostalCodeMatches, types.CityPostalCode{
								City:                cm.Value,
								PostalCode:          pc,
								CountryCode:         "si",
								CitySimilarity:      cm.Similarity,
								WasPartialCityMatch: cm.WasPartial,
							})
						}
					}
				}
			}
			rows.Close()
		}

		if len(result.CityPostalCodeMatches) > 0 {
			result.CityPostalCodeMatch = true
		}
		reChan <- result
	}()
	return reChan
}

func like(s string) string {
	// escape % and _ minimally; SQLite uses \\ for ESCAPE not by default; keep simple
	s = strings.ToLower(s)
	return fmt.Sprintf("%%%s%%", s)
}
