package normalization

import (
	"fmt"
	"regexp"
	"strings"
)

type DK struct {
	reCompanySuffix       *regexp.Regexp // can contain c/o information and should be ignored
	reBuildingAfterNumber *regexp.Regexp // can contain additional location info, such as "entrance", "building" etc.
	reFloor               *regexp.Regexp
	reDoor                *regexp.Regexp
	reCityGeo             *regexp.Regexp // Bigger cities can be suffixed with a geographical location description, e.g. København NV (Copenhagen North West)
	reCrap                *regexp.Regexp
	rePostalCode          *regexp.Regexp
}

func newDK() (*DK, error) {
	reCompanySuffix, errCompanySuffix := regexp.Compile("\\s*( v/| c/o| att)")
	if errCompanySuffix != nil {
		return nil, errCompanySuffix
	}
	reBuilding, errBuildings := regexp.Compile("\\W(sal|etage|floor|kl|kld|kælder|st|stuen|parterre|dør|port|opg|opgang|indgang|bygning|tv|th|bygn|mf)\\b")
	if errBuildings != nil {
		return nil, errBuildings
	}

	reFloor, errFloor := regexp.Compile("(\\W)\\b(kld|kl|kælder\\w*|st|stuen|parterren)\\b")
	if errFloor != nil {
		return nil, errFloor
	}

	reDoor, errDoor := regexp.Compile("(\\d.*)(?:\\btv|\\bt\\W+v|v\\.|\\bv|\\bth|\\bt\\W+h|h\\.|\\bh|mf)\\b")
	if errDoor != nil {
		return nil, errDoor
	}

	reCityGeo, errGeo := regexp.Compile("\\W(sv|nv|sø|nø|s|n|v|ø)")
	if errGeo != nil {
		return nil, errGeo
	}

	rePostalCode, errPostalCode := regexp.Compile("[+/(){}\\[\\]<>!§$%&=?*#€¿_\",:;a-zA-Z- ]")
	if errPostalCode != nil {
		return nil, errPostalCode
	}

	reCrap, errCrap := regexp.Compile("[+/(){}\\[\\]<>!§'$%&=?*#€¿_\":;0-9]")
	if errCrap != nil {
		return nil, errCrap
	}

	return &DK{
		reCompanySuffix:       reCompanySuffix,
		reBuildingAfterNumber: reBuilding,
		reFloor:               reFloor,
		reDoor:                reDoor,
		reCityGeo:             reCityGeo,
		rePostalCode:          rePostalCode,
		reCrap:                reCrap,
	}, nil
}

func (dk *DK) GetCountryCode() string {
	return "dk"
}

func (dk *DK) City(s string) (string, error) {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "oe", "ø")
	s = strings.ReplaceAll(s, "ae", "æ")
	s = dk.reCityGeo.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)

	return s, nil
}

func (dk *DK) PostalCode(s string) (string, error) {
	s = dk.rePostalCode.ReplaceAllString(s, "")

	if len(s) != 4 {
		if len(s) > 4 {
			return s[:4], fmt.Errorf("not valid postalcode")
		}
		return s, fmt.Errorf("not valid postalcode")
	}

	return s, nil
}

func (dk *DK) Street(s string) ([]string, error) {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "\n", "|")
	s = strings.ReplaceAll(s, ",", "|")

	s = strings.ReplaceAll(s, "oe", "ø")
	s = strings.ReplaceAll(s, "ae", "æ")

	s = dk.reDoor.ReplaceAllString(s, " ")
	s = dk.reCompanySuffix.ReplaceAllString(s, " ")
	s = dk.reBuildingAfterNumber.ReplaceAllString(s, " ")
	s = dk.reFloor.ReplaceAllString(s, " ")

	s = dk.reCrap.ReplaceAllString(s, " ")

	addressSlice := strings.Split(s, "|")
	var cleanAddressSlice []string
	for _, v := range addressSlice {
		var cleanedParts []string
		addressPartSlices := strings.Split(v, " ")
		for _, p := range addressPartSlices {
			if len(p) > 1 {
				cleanedParts = append(cleanedParts, p)
			}
		}
		cleanPart := strings.Join(cleanedParts, " ")
		if len(cleanPart) > 1 {
			cleanAddressSlice = append(cleanAddressSlice, cleanPart)
		}
	}

	return cleanAddressSlice, nil
}
