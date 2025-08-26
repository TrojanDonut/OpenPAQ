#!/bin/bash

# Test script for OpenPAQ address validation
# Tests both Slovenian addresses (SQLite) and global addresses (Nominatim)

BASE_URL="http://127.0.0.1:8080"
API_ENDPOINT="$BASE_URL/api/v1/check"

echo "=== OpenPAQ Address Validation Test Suite ==="
echo "Base URL: $BASE_URL"
echo ""

# Function to test an address
test_address() {
    local street="$1"
    local city="$2"
    local postal_code="$3"
    local country_code="$4"
    local description="$5"
    
    echo "Testing: $description"
    echo "  Street: $street"
    echo "  City: $city"
    echo "  Postal Code: $postal_code"
    echo "  Country: $country_code"
    
    # URL encode the parameters
    street_encoded=$(echo "$street" | sed 's/ /%20/g')
    city_encoded=$(echo "$city" | sed 's/ /%20/g')
    
    # Make the API call
    response=$(curl -s "$API_ENDPOINT?street=$street_encoded&city=$city_encoded&postal_code=$postal_code&country_code=$country_code")
    
    # Check if the request was successful
    if [ $? -eq 0 ]; then
        echo "  Response: $response"
        
        # Extract key fields using jq if available
        if command -v jq &> /dev/null; then
            street_matched=$(echo "$response" | jq -r '.street_matched // "unknown"')
            city_matched=$(echo "$response" | jq -r '.city_matched // "unknown"')
            postal_matched=$(echo "$response" | jq -r '.postal_code_matched // "unknown"')
            country_matched=$(echo "$response" | jq -r '.country_code_matched // "unknown"')
            
            echo "  Results:"
            echo "    Street matched: $street_matched"
            echo "    City matched: $city_matched"
            echo "    Postal code matched: $postal_matched"
            echo "    Country code matched: $country_matched"
        fi
    else
        echo "  Error: Failed to make API request"
    fi
    
    echo ""
}

# Function to test with debug details
test_address_debug() {
    local street="$1"
    local city="$2"
    local postal_code="$3"
    local country_code="$4"
    local description="$5"
    
    echo "Testing (with debug): $description"
    
    # URL encode the parameters
    street_encoded=$(echo "$street" | sed 's/ /%20/g')
    city_encoded=$(echo "$city" | sed 's/ /%20/g')
    
    # Make the API call with debug_details=true
    response=$(curl -s "$API_ENDPOINT?street=$street_encoded&city=$city_encoded&postal_code=$postal_code&country_code=$country_code&debug_details=true")
    
    if [ $? -eq 0 ]; then
        echo "  Debug Response: $response"
    else
        echo "  Error: Failed to make API request"
    fi
    
    echo ""
}

echo "=== Testing Slovenian Addresses (SQLite Database) ==="
echo "These tests use the local SQLite database for fast, reliable validation."
echo ""

# Test 1: Valid Slovenian address from our CSV data
test_address "Trzinska cesta" "Mengeš" "1234" "si" "Valid Slovenian address (Mengeš)"

# Test 2: Another valid Slovenian address
test_address "Cesta VIII" "Grič" "1310" "si" "Valid Slovenian address (Grič)"

# Test 3: Slovenian address with partial postal code
test_address "Šentviška pot" "Čatež ob Savi" "8250" "si" "Valid Slovenian address (Čatež ob Savi)"

# Test 4: Invalid Slovenian address
test_address "Nonexistent Street" "Nonexistent City" "9999" "si" "Invalid Slovenian address"

echo "=== Testing Global Addresses (Nominatim) ==="
echo "These tests use the Nominatim API for global address validation."
echo ""

# Test 5: German address
test_address "Unter den Linden" "Berlin" "10117" "de" "Valid German address (Berlin)"

# Test 6: French address
test_address "Champs-Élysées" "Paris" "75008" "fr" "Valid French address (Paris)"

# Test 7: UK address
test_address "Oxford Street" "London" "W1C 1AP" "gb" "Valid UK address (London)"

# Test 8: US address
test_address "Times Square" "New York" "10036" "us" "Valid US address (New York)"

# Test 9: Invalid global address
test_address "Fake Street 123" "Fake City" "00000" "us" "Invalid global address"

echo "=== Testing Edge Cases ==="
echo ""

# Test 10: Empty fields
test_address "" "" "" "si" "Empty fields (Slovenian)"

# Test 11: Very long street name
test_address "This is a very long street name that might cause issues with the API" "Ljubljana" "1000" "si" "Long street name"

# Test 12: Special characters
test_address "Študentovska ulica" "Ljubljana" "1000" "si" "Slovenian address with special characters"

echo "=== Debug Testing ==="
echo "Testing with debug_details=true to see detailed matching information"
echo ""

# Test with debug details for a valid Slovenian address
test_address_debug "Trzinska cesta" "Mengeš" "1234" "si" "Valid Slovenian address with debug details"

# Test with debug details for a valid global address
test_address_debug "Unter den Linden" "Berlin" "10117" "de" "Valid German address with debug details"

echo "=== Performance Testing ==="
echo "Testing response times for different address types"
echo ""

# Test Slovenian address performance
echo "Testing Slovenian address response time:"
time curl -s "$API_ENDPOINT?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si" > /dev/null

echo ""
echo "Testing global address response time:"
time curl -s "$API_ENDPOINT?street=Unter%20den%20Linden&city=Berlin&postal_code=10117&country_code=de" > /dev/null