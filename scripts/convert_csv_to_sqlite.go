package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SlovenianAddress struct {
	FeatureID              string
	EIDNaslov              string
	ObcinaSifra            string
	ObcinaNaziv            string
	ObcinaNazivDJ          string
	NaseljeSifra           string
	NaseljeNaziv           string
	NaseljeNazivDJ         string
	UlicaSifra             string
	UlicaNaziv             string
	UlicaNazivDJ           string
	PostniOkolisSifra      string
	PostniOkolisNaziv      string
	PostniOkolisNazivDJ    string
	HSStevilka             string
	HSDodatek              string
	STStanovanja           string
	E                      string
	N                      string
	EIDObcina              string
	EIDNaselje             string
	EIDUlica               string
	EIDPostniOkolis        string
	EIDHisnaStevilka       string
	EIDStanovanje          string
	EIDStavba              string
	EIDCetrtnaSkupnost     string
	EIDDzVolisce           string
	EIDKrajevnaSkupnost    string
	EIDLokalnoVolisce      string
	EIDLokalnaVolilnaEnota string
	EIDSolskiOkolis        string
	EIDStatisticnaRegija   string
	EIDUpravnaEnota        string
	EIDVaskaSkupnost       string
	EIDVolilnaEnotaDZ      string
	EIDVolilniOkraj        string
	EIDKohezijskaRegija    string
	DatumSYS               string
}

func main() {
	// Open CSV file
	csvFile, err := os.Open("RN_SLO_NASLOVI_register_naslovov_20250817.csv")
	if err != nil {
		log.Fatal("Error opening CSV file:", err)
	}
	defer csvFile.Close()

	// Create SQLite database
	db, err := sql.Open("sqlite3", "slovenian_addresses.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Create table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS slovenian_addresses (
		feature_id TEXT PRIMARY KEY,
		eid_naslov TEXT,
		obcina_sifra TEXT,
		obcina_naziv TEXT,
		obcina_naziv_dj TEXT,
		naselje_sifra TEXT,
		naselje_naziv TEXT,
		naselje_naziv_dj TEXT,
		ulica_sifra TEXT,
		ulica_naziv TEXT,
		ulica_naziv_dj TEXT,
		postni_okolis_sifra TEXT,
		postni_okolis_naziv TEXT,
		postni_okolis_naziv_dj TEXT,
		hs_stevilka TEXT,
		hs_dodatek TEXT,
		st_stanovanja TEXT,
		e TEXT,
		n TEXT,
		eid_obcina TEXT,
		eid_naselje TEXT,
		eid_ulica TEXT,
		eid_postni_okolis TEXT,
		eid_hisna_stevilka TEXT,
		eid_stanovanje TEXT,
		eid_stavba TEXT,
		eid_cetrtna_skupnost TEXT,
		eid_dz_volisce TEXT,
		eid_krajevna_skupnost TEXT,
		eid_lokalno_volisce TEXT,
		eid_lokalna_volilna_enota TEXT,
		eid_solski_okolis TEXT,
		eid_statisticna_regija TEXT,
		eid_upravna_enota TEXT,
		eid_vaska_skupnost TEXT,
		eid_volilna_enota_dz TEXT,
		eid_volilni_okraj TEXT,
		eid_kohezijska_regija TEXT,
		datum_sys TEXT
	);
	
	CREATE INDEX IF NOT EXISTS idx_obcina_naziv ON slovenian_addresses(obcina_naziv);
	CREATE INDEX IF NOT EXISTS idx_naselje_naziv ON slovenian_addresses(naselje_naziv);
	CREATE INDEX IF NOT EXISTS idx_ulica_naziv ON slovenian_addresses(ulica_naziv);
	CREATE INDEX IF NOT EXISTS idx_postni_okolis_naziv ON slovenian_addresses(postni_okolis_naziv);
	CREATE INDEX IF NOT EXISTS idx_postni_okolis_sifra ON slovenian_addresses(postni_okolis_sifra);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	// Read CSV
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Skip header
	_, err = reader.Read()
	if err != nil {
		log.Fatal("Error reading header:", err)
	}

	// Prepare insert statement
	insertSQL := `
	INSERT OR REPLACE INTO slovenian_addresses (
		feature_id, eid_naslov, obcina_sifra, obcina_naziv, obcina_naziv_dj,
		naselje_sifra, naselje_naziv, naselje_naziv_dj, ulica_sifra, ulica_naziv, ulica_naziv_dj,
		postni_okolis_sifra, postni_okolis_naziv, postni_okolis_naziv_dj, hs_stevilka, hs_dodatek,
		st_stanovanja, e, n, eid_obcina, eid_naselje, eid_ulica, eid_postni_okolis,
		eid_hisna_stevilka, eid_stanovanje, eid_stavba, eid_cetrtna_skupnost, eid_dz_volisce,
		eid_krajevna_skupnost, eid_lokalno_volisce, eid_lokalna_volilna_enota, eid_solski_okolis,
		eid_statisticna_regija, eid_upravna_enota, eid_vaska_skupnost, eid_volilna_enota_dz,
		eid_volilni_okraj, eid_kohezijska_regija, datum_sys
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	// Read and insert records
	recordCount := 0
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		// Ensure we have enough fields
		if len(record) < 39 {
			continue // Skip malformed records
		}

		// Clean and normalize data
		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		_, err = stmt.Exec(
			record[0], record[1], record[2], record[3], record[4],
			record[5], record[6], record[7], record[8], record[9], record[10],
			record[11], record[12], record[13], record[14], record[15],
			record[16], record[17], record[18], record[19], record[20], record[21], record[22],
			record[23], record[24], record[25], record[26], record[27],
			record[28], record[29], record[30], record[31], record[32],
			record[33], record[34], record[35], record[36], record[37], record[38],
		)

		if err != nil {
			log.Printf("Error inserting record %d: %v", recordCount, err)
			continue
		}

		recordCount++
		if recordCount%10000 == 0 {
			fmt.Printf("Processed %d records\n", recordCount)
		}
	}

	fmt.Printf("Successfully imported %d records to SQLite database\n", recordCount)

	// Print some statistics
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM slovenian_addresses").Scan(&count)
	if err != nil {
		log.Printf("Error counting records: %v", err)
	} else {
		fmt.Printf("Total records in database: %d\n", count)
	}

	// Print unique cities count
	var cityCount int
	err = db.QueryRow("SELECT COUNT(DISTINCT naselje_naziv) FROM slovenian_addresses WHERE naselje_naziv != ''").Scan(&cityCount)
	if err != nil {
		log.Printf("Error counting cities: %v", err)
	} else {
		fmt.Printf("Unique cities: %d\n", cityCount)
	}

	// Print unique streets count
	var streetCount int
	err = db.QueryRow("SELECT COUNT(DISTINCT ulica_naziv) FROM slovenian_addresses WHERE ulica_naziv != ''").Scan(&streetCount)
	if err != nil {
		log.Printf("Error counting streets: %v", err)
	} else {
		fmt.Printf("Unique streets: %d\n", streetCount)
	}

	// Print unique postal codes count
	var postalCount int
	err = db.QueryRow("SELECT COUNT(DISTINCT postni_okolis_sifra) FROM slovenian_addresses WHERE postni_okolis_sifra != ''").Scan(&postalCount)
	if err != nil {
		log.Printf("Error counting postal codes: %v", err)
	} else {
		fmt.Printf("Unique postal codes: %d\n", postalCount)
	}
}
