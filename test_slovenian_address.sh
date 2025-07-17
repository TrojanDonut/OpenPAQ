#!/bin/bash

# OpenPAQ Slovenian Address Tester
# Usage: ./test_slovenian_address.sh "Tehnološki park 21, 1000 Ljubljana"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if address is provided
if [ $# -eq 0 ]; then
    echo -e "${RED}Error: No address provided${NC}"
    echo "Usage: $0 \"<address>\""
    echo "Example: $0 \"Tehnološki park 21, 1000 Ljubljana\""
    exit 1
fi

# Get the full address
FULL_ADDRESS="$1"

echo -e "${BLUE}=== OpenPAQ Slovenian Address Validator ===${NC}"
echo -e "${YELLOW}Input address:${NC} $FULL_ADDRESS"
echo

# Split by comma first
IFS=',' read -r STREET_PART POSTAL_CITY <<< "$FULL_ADDRESS"

if [ -z "$POSTAL_CITY" ]; then
    echo -e "${RED}Error: Could not parse address format${NC}"
    echo "Expected format: \"Street [Number], PostalCode City\""
    echo "Example: \"Tehnološki park 21, 1000 Ljubljana\""
    exit 1
fi

# Extract postal code and city from the second part
if [[ $POSTAL_CITY =~ ^[[:space:]]*([0-9]{4})[[:space:]]+(.+)$ ]]; then
    POSTAL_CODE="${BASH_REMATCH[1]}"
    CITY="${BASH_REMATCH[2]}"
    STREET="$STREET_PART"
else
    echo -e "${RED}Error: Could not parse postal code and city${NC}"
    echo "Expected format: \"Street [Number], PostalCode City\""
    echo "Example: \"Tehnološki park 21, 1000 Ljubljana\""
    exit 1
fi

echo -e "${YELLOW}Parsed components:${NC}"
echo "  Street: $STREET"
echo "  Postal Code: $POSTAL_CODE"
echo "  City: $CITY"
echo "  Country: SI"
echo

# URL encode the components
STREET_ENCODED=$(echo "$STREET" | sed 's/ /+/g')
CITY_ENCODED=$(echo "$CITY" | sed 's/ /+/g')

# Make the API call
echo -e "${YELLOW}Making API call to OpenPAQ...${NC}"
echo

API_URL="http://127.0.0.1:8001/api/v1/check?street=${STREET_ENCODED}&postal_code=${POSTAL_CODE}&city=${CITY_ENCODED}&country_code=SI&debug_details=false"

# Get the response
RESPONSE=$(curl -s "$API_URL")

# Check if curl was successful
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to connect to OpenPAQ API${NC}"
    echo "Make sure OpenPAQ is running on http://127.0.0.1:8001"
    exit 1
fi

# Parse JSON response (simple parsing)
STREET_MATCHED=$(echo "$RESPONSE" | grep -o '"street_matched":[^,]*' | cut -d':' -f2)
CITY_MATCHED=$(echo "$RESPONSE" | grep -o '"city_matched":[^,]*' | cut -d':' -f2)
POSTAL_MATCHED=$(echo "$RESPONSE" | grep -o '"postal_code_matched":[^,]*' | cut -d':' -f2)
CITY_POSTAL_MATCHED=$(echo "$RESPONSE" | grep -o '"city_to_postal_code_matched":[^,]*' | cut -d':' -f2)
COUNTRY_MATCHED=$(echo "$RESPONSE" | grep -o '"country_code_matched":[^,]*' | cut -d':' -f2)

echo -e "${BLUE}=== Validation Results ===${NC}"

# Display results with colors
if [ "$STREET_MATCHED" = "true" ]; then
    echo -e "  Street: ${GREEN}✓ Valid${NC}"
else
    echo -e "  Street: ${RED}✗ Invalid${NC}"
fi

if [ "$CITY_MATCHED" = "true" ]; then
    echo -e "  City: ${GREEN}✓ Valid${NC}"
else
    echo -e "  City: ${RED}✗ Invalid${NC}"
fi

if [ "$POSTAL_MATCHED" = "true" ]; then
    echo -e "  Postal Code: ${GREEN}✓ Valid${NC}"
else
    echo -e "  Postal Code: ${RED}✗ Invalid${NC}"
fi

if [ "$CITY_POSTAL_MATCHED" = "true" ]; then
    echo -e "  City-Postal Match: ${GREEN}✓ Valid${NC}"
else
    echo -e "  City-Postal Match: ${RED}✗ Invalid${NC}"
fi

if [ "$COUNTRY_MATCHED" = "true" ]; then
    echo -e "  Country: ${GREEN}✓ Valid${NC}"
else
    echo -e "  Country: ${RED}✗ Invalid${NC}"
fi

echo
echo -e "${BLUE}=== Raw API Response ===${NC}"
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE" 