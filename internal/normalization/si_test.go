package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCodeSI(t *testing.T) {
	siNormalizer, _ := NewSI()

	countryCode := siNormalizer.GetCountryCode()

	if countryCode != "si" {
		t.Errorf("got %s, want si", countryCode)
	}
}

func TestPostalCodeSI(t *testing.T) {
	siNormalizer, _ := NewSI()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Valid 4-digit postal code", input: "1000", expect: "1000"},
		{name: "Postal code with letters should be removed", input: "1000Ljubljana", expect: "1000"},
		{name: "Postal code with special characters should be removed", input: "1000-123", expect: "1000"},
		{name: "Postal code with spaces should be removed", input: "1000 123", expect: "1000"},
		{name: "Postal code with leading zeros", input: "0100", expect: "0100"},
		{name: "Postal code with mixed characters", input: "SI-1000", expect: "1000"},
		{name: "Postal code with dots", input: "1000.123", expect: "1000"},
		{name: "Postal code with parentheses", input: "(1000)", expect: "1000"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := siNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeExpectErrorSI(t *testing.T) {
	siNormalizer, _ := NewSI()

	tests := []struct {
		name        string
		input       string
		expect      string
		expectError bool
	}{
		{name: "Postal code too long", input: "10000", expect: "1000", expectError: false},
		{name: "Postal code too short", input: "100", expect: "100", expectError: true},
		{name: "Empty postal code", input: "", expect: "", expectError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := siNormalizer.PostalCode(test.input)
			if (err != nil) != test.expectError {
				t.Errorf("got error %v, expectError %v", err, test.expectError)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestCitySI(t *testing.T) {
	siNormalizer, _ := NewSI()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be lower case", input: "Ljubljana", expect: "ljubljana"},
		{name: "Whitespaces should be trimmed", input: "   Ljubljana   ", expect: "ljubljana"},
		{name: "Special characters should be removed", input: "Ljubljana+/(){}[]<>!§'$%&=?*#€¿_\":;12345", expect: "ljubljana"},
		{name: "Slovenian diacritics should be handled", input: "Črnomelj", expect: "crnomelj"},
		{name: "Multiple Slovenian diacritics", input: "Škofja Loka", expect: "skofja loka"},
		{name: "Numbers should be removed", input: "Ljubljana 1000", expect: "ljubljana"},
		{name: "Mixed case with diacritics", input: "Koper-Portorož", expect: "koper-portoroz"},
		{name: "City with common suffixes", input: "Maribor mesto", expect: "maribor mesto"},
		{name: "City with abbreviations", input: "Celje ob Savinji", expect: "celje ob savinji"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := siNormalizer.City(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreetSI(t *testing.T) {
	siNormalizer, _ := NewSI()

	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{name: "Basic street name", input: "Slovenska cesta", expect: []string{"slovenska cesta"}},
		{name: "Street with abbreviation", input: "Slovenska ul.", expect: []string{"slovenska ulica"}},
		{name: "Street with number", input: "Slovenska cesta 1", expect: []string{"slovenska cesta"}},
		{name: "Split at new line", input: "Slovenska cesta\nLjubljana", expect: []string{"slovenska cesta", "ljubljana"}},
		{name: "Split at comma", input: "Slovenska cesta,Ljubljana", expect: []string{"slovenska cesta", "ljubljana"}},
		{name: "Special characters should be removed", input: "Slovenska cesta+/(){}[]<>!§'$%&=?*#€¿_\":;12345", expect: []string{"slovenska cesta"}},
		{name: "Slovenian diacritics should be handled", input: "Čopova ulica", expect: []string{"copova ulica"}},
		{name: "Multiple street types", input: "Trg republike", expect: []string{"trg republike"}},
		{name: "Street with building info", input: "Slovenska cesta 1, stavba A", expect: []string{"slovenska cesta", "stavba"}},
		{name: "Street with floor info", input: "Slovenska cesta 1, 2. nadstropje", expect: []string{"slovenska cesta", "nadstropje"}},
		{name: "Street abbreviations", input: "Slovenska ul. 15", expect: []string{"slovenska ulica"}},
		{name: "Street with cesta abbreviation", input: "Slovenska c. 15", expect: []string{"slovenska cesta"}},
		{name: "Street with trg abbreviation", input: "Trg republike 1", expect: []string{"trg republike"}},
		{name: "Complex address with multiple parts", input: "Slovenska cesta 1, Ljubljana, 1000", expect: []string{"slovenska cesta", "ljubljana"}},
		{name: "Street with apartment info", input: "Slovenska cesta 1, stanovanje 5", expect: []string{"slovenska cesta", "stanovanje"}},
		{name: "Street with entrance info", input: "Slovenska cesta 1, vhod A", expect: []string{"slovenska cesta", "vhod"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := siNormalizer.Street(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if !slices.Equal(test.expect, res) {
				t.Errorf("got %v, want %v", res, test.expect)
			}
		})
	}
}

// Benchmark tests for performance comparison
func BenchmarkPostalCodeSI(b *testing.B) {
	siNormalizer, _ := NewSI()
	for i := 0; i < b.N; i++ {
		_, _ = siNormalizer.PostalCode("1000")
	}
}

func BenchmarkCitySI(b *testing.B) {
	siNormalizer, _ := NewSI()
	for i := 0; i < b.N; i++ {
		_, _ = siNormalizer.City("Ljubljana")
	}
}

func BenchmarkStreetSI(b *testing.B) {
	siNormalizer, _ := NewSI()
	for i := 0; i < b.N; i++ {
		_, _ = siNormalizer.Street("Slovenska cesta 1")
	}
}
