package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCode(t *testing.T) {
	dkNormalizer, _ := newDK()

	countryCode := dkNormalizer.GetCountryCode()

	if countryCode != "dk" {
		t.Errorf("got %s, want dk", countryCode)
	}
}

func TestCity(t *testing.T) {
	dkNormalizer, _ := newDK()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Want only small letters", input: "ACityWithCapitalLetters", expect: "acitywithcapitalletters"},
		{name: "oe should be replaced", input: "ACityWithOeWithin", expect: "acitywithøwithin"},
		{name: "ae should be replaced", input: "ACityWithAeWithin", expect: "acitywithæwithin"},
		{name: "sv should be removed at the end", input: "City SV", expect: "city"},
		{name: "sø should be removed at the end", input: "City sø", expect: "city"},
		{name: "nø should be removed at the end", input: "City nø", expect: "city"},
		{name: "s should be removed at the end", input: "City S", expect: "city"},
		{name: "N should be removed at the end", input: "City N", expect: "city"},
		{name: "V should be removed at the end", input: "City V", expect: "city"},
		{name: "ø should be removed at the end", input: "City ø", expect: "city"},
		{name: "sv should not be removed at the beginning", input: "SV City", expect: "sv city"},
		{name: "sv should not be removed at the beginning", input: "SVCity", expect: "svcity"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.City(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeNoError(t *testing.T) {
	dkNormalizer, _ := newDK()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be removed", input: "1234abcdefghijklmnopqrstuvwxyz", expect: "1234"},
		{name: "Capital letters should be removed", input: "1234ABCDEFGHIJKLMNOPQRSTUVWXYZ", expect: "1234"},
		{name: "Special characters should be removed", input: "1234+/(){}[]<>!§$%&=?*#€¿_\",:;-", expect: "1234"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeExpectError(t *testing.T) {
	dkNormalizer, _ := newDK()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Postal code to long", input: "12345", expect: "1234"},
		{name: "Postal code to short", input: "123", expect: "123"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.PostalCode(test.input)
			if err == nil {
				t.Error("expect error got none")
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreet(t *testing.T) {
	dkNormalizer, _ := newDK()

	tests := []struct {
		name            string
		input           string
		expectedContent []string
	}{
		{name: "no separator", input: "JustAStreet", expectedContent: []string{"justastreet"}},
		{name: "spaces", input: "Just A Street", expectedContent: []string{"just street"}},
		{name: "comma", input: "Just , Street", expectedContent: []string{"just", "street"}},
		{name: "new line", input: "Just \n Street", expectedContent: []string{"just", "street"}},
		{name: "oe", input: "Just Streetoe", expectedContent: []string{"just streetø"}},
		{name: "ae", input: "Just Streetae", expectedContent: []string{"just streetæ"}},
		{name: "remove v at the end", input: "Just Street v", expectedContent: []string{"just street"}},
		{name: "remove v at the end of word", input: "Just Streetv", expectedContent: []string{"just streetv"}},
		{name: "c/o at the end", input: "Just Street c/o", expectedContent: []string{"just street"}},
		{name: "att at the end", input: "Just Street att", expectedContent: []string{"just street"}},
		{name: "att at the end of word", input: "Just Streetatt", expectedContent: []string{"just streetatt"}},
		{name: "att at the end of word", input: "Just Streetatt", expectedContent: []string{"just streetatt"}},
		{name: "building context", input: "Just Street sal etage floor kl kld kælder st stuen parterre dør port opg opgang indgang bygning bygn tv th mf", expectedContent: []string{"just street"}},
		{name: "door context test-1", input: "Test Vænge 9, opgang 4, 2. tv blub", expectedContent: []string{"test vænge blub"}},
		{name: "door context test-2", input: "lurchschlumpf 2. tv asdf", expectedContent: []string{"lurchschlumpf asdf"}},
		{name: "door context test-3", input: "lurchschlumpf 1.tv.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-4", input: "lurchschlumpf 1.t.v.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-5", input: "lurchschlumpf 1.v.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-6", input: "lurchschlumpf 1. v.", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-7", input: "lurchschlumpf 1. th", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-8", input: "lurchschlumpf 1. h", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-9", input: "lurchschlumpf 1 h", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-10", input: "lurchschlumpf 1 v", expectedContent: []string{"lurchschlumpf"}},
		{name: "door context test-11", input: "lurchschlumpf 1. mf", expectedContent: []string{"lurchschlumpf"}},
		{name: "floor context", input: "lurchschlumpf kld kl kælder st stuen parterren", expectedContent: []string{"lurchschlumpf"}},
		{name: "floor context", input: "lurchschlumpfkld", expectedContent: []string{"lurchschlumpfkld"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := dkNormalizer.Street(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if !slices.Equal(test.expectedContent, res) {
				t.Errorf("got %v, want %v", res, test.expectedContent)
			}
		})
	}
}
