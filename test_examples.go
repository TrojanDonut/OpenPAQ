package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// AddressTest represents a test case
type AddressTest struct {
	Street     string
	City       string
	PostalCode string
	CountryCode string
	Description string
	Expected   bool // Expected to be valid
}

// APIResponse represents the API response structure
type APIResponse struct {
	Street                string `json:"street"`
	City                  string `json:"city"`
	PostalCode            string `json:"postal_code"`
	CountryCode           string `json:"country_code"`
	StreetMatched         bool   `json:"street_matched"`
	CityMatched           bool   `json:"city_matched"`
	PostalCodeMatched     bool   `json:"postal_code_matched"`
	CityToPostalCodeMatch bool   `json:"city_to_postal_code_matched"`
	CountryCodeMatched    bool   `json:"country_code_matched"`
	Version               string `json:"version"`
}

func main() {
	baseURL := "http://127.0.0.1:8080"
	
	// Test cases for Slovenian addresses (SQLite database)
	slovenianTests := []AddressTest{
		{
			Street:      "Trzinska cesta",
			City:        "Mengeš",
			PostalCode:  "1234",
			CountryCode: "si",
			Description: "Valid Slovenian address (Mengeš)",
			Expected:    true,
		},
		{
			Street:      "Cesta VIII",
			City:        "Grič",
			PostalCode:  "1310",
			CountryCode: "si",
			Description: "Valid Slovenian address (Grič)",
			Expected:    true,
		},
		{
			Street:      "Šentviška pot",
			City:        "Čatež ob Savi",
			PostalCode:  "8250",
			CountryCode: "si",
			Description: "Valid Slovenian address (Čatež ob Savi)",
			Expected:    true,
		},
		{
			Street:      "Nonexistent Street",
			City:        "Nonexistent City",
			PostalCode:  "9999",
			CountryCode: "si",
			Description: "Invalid Slovenian address",
			Expected:    false,
		},
	}

	// Test cases for global addresses (Nominatim)
	globalTests := []AddressTest{
		{
			Street:      "Unter den Linden",
			City:        "Berlin",
			PostalCode:  "10117",
			CountryCode: "de",
			Description: "Valid German address (Berlin)",
			Expected:    true,
		},
		{
			Street:      "Champs-Élysées",
			City:        "Paris",
			PostalCode:  "75008",
			CountryCode: "fr",
			Description: "Valid French address (Paris)",
			Expected:    true,
		},
		{
			Street:      "Oxford Street",
			City:        "London",
			PostalCode:  "W1C 1AP",
			CountryCode: "gb",
			Description: "Valid UK address (London)",
			Expected:    true,
		},
		{
			Street:      "Times Square",
			City:        "New York",
			PostalCode:  "10036",
			CountryCode: "us",
			Description: "Valid US address (New York)",
			Expected:    true,
		},
		{
			Street:      "Fake Street 123",
			City:        "Fake City",
			PostalCode:  "00000",
			CountryCode: "us",
			Description: "Invalid global address",
			Expected:    false,
		},
	}

	fmt.Println("=== OpenPAQ Address Validation Test Examples ===")
	fmt.Println()

	// Test Slovenian addresses
	fmt.Println("=== Testing Slovenian Addresses (SQLite Database) ===")
	fmt.Println("These tests use the local SQLite database for fast, reliable validation.")
	fmt.Println()
	
	for i, test := range slovenianTests {
		fmt.Printf("Test %d: %s\n", i+1, test.Description)
		result := testAddress(baseURL, test)
		printResult(result, test)
		fmt.Println()
	}

	// Test global addresses
	fmt.Println("=== Testing Global Addresses (Nominatim) ===")
	fmt.Println("These tests use the Nominatim API for global address validation.")
	fmt.Println()
	
	for i, test := range globalTests {
		fmt.Printf("Test %d: %s\n", i+1, test.Description)
		result := testAddress(baseURL, test)
		printResult(result, test)
		fmt.Println()
	}

	// Performance comparison
	fmt.Println("=== Performance Comparison ===")
	fmt.Println()
	
	// Test Slovenian performance
	start := time.Now()
	testAddress(baseURL, slovenianTests[0])
	slovenianTime := time.Since(start)
	
	// Test global performance
	start = time.Now()
	testAddress(baseURL, globalTests[0])
	globalTime := time.Since(start)
	
	fmt.Printf("Slovenian address (SQLite): %v\n", slovenianTime)
	fmt.Printf("Global address (Nominatim): %v\n", globalTime)
	fmt.Printf("Performance improvement: %.2fx faster\n", float64(globalTime)/float64(slovenianTime))
}

func testAddress(baseURL string, test AddressTest) *APIResponse {
	// Build URL with parameters
	params := url.Values{}
	params.Add("street", test.Street)
	params.Add("city", test.City)
	params.Add("postal_code", test.PostalCode)
	params.Add("country_code", test.CountryCode)
	
	apiURL := fmt.Sprintf("%s/api/v1/check?%s", baseURL, params.Encode())
	
	// Make HTTP request
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("  Error making request: %v\n", err)
		return nil
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("  Error reading response: %v\n", err)
		return nil
	}
	
	// Parse JSON response
	var result APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("  Error parsing JSON: %v\n", err)
		return nil
	}
	
	return &result
}

func printResult(result *APIResponse, test AddressTest) {
	if result == nil {
		fmt.Println("  ❌ Failed to get result")
		return
	}
	
	fmt.Printf("  Street: %s\n", result.Street)
	fmt.Printf("  City: %s\n", result.City)
	fmt.Printf("  Postal Code: %s\n", result.PostalCode)
	fmt.Printf("  Country Code: %s\n", result.CountryCode)
	fmt.Println("  Results:")
	fmt.Printf("    Street matched: %t\n", result.StreetMatched)
	fmt.Printf("    City matched: %t\n", result.CityMatched)
	fmt.Printf("    Postal code matched: %t\n", result.PostalCodeMatched)
	fmt.Printf("    City-to-postal match: %t\n", result.CityToPostalCodeMatch)
	fmt.Printf("    Country code matched: %t\n", result.CountryCodeMatched)
	
	// Determine if the address is valid (all matches true)
	isValid := result.StreetMatched && result.CityMatched && result.PostalCodeMatched && 
			   result.CityToPostalCodeMatch && result.CountryCodeMatched
	
	if isValid == test.Expected {
		fmt.Printf("  ✅ Expected: %t, Got: %t (PASS)\n", test.Expected, isValid)
	} else {
		fmt.Printf("  ❌ Expected: %t, Got: %t (FAIL)\n", test.Expected, isValid)
	}
}

// Example usage functions
func ExampleSlovenianAddress() {
	// Example: Test a Slovenian address
	test := AddressTest{
		Street:      "Trzinska cesta",
		City:        "Mengeš",
		PostalCode:  "1234",
		CountryCode: "si",
		Description: "Example Slovenian address",
		Expected:    true,
	}
	
	result := testAddress("http://127.0.0.1:8080", test)
	if result != nil {
		fmt.Printf("Slovenian address validation result: %+v\n", result)
	}
}

func ExampleGlobalAddress() {
	// Example: Test a global address
	test := AddressTest{
		Street:      "Unter den Linden",
		City:        "Berlin",
		PostalCode:  "10117",
		CountryCode: "de",
		Description: "Example German address",
		Expected:    true,
	}
	
	result := testAddress("http://127.0.0.1:8080", test)
	if result != nil {
		fmt.Printf("Global address validation result: %+v\n", result)
	}
}
