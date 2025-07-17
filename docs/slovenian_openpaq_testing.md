# Slovenian Address Validation with OpenPAQ

## 1. Overview
This document describes how we tested OpenPAQ for Slovenian address validation, including:
- How we set up and ran OpenPAQ
- How we implemented Slovenian normalization
- How we ran and fixed the tests
- Example API usage and results

---

## 2. Running OpenPAQ

We started the OpenPAQ server locally using the following command:

```bash
CACHE_ENABLED=false \
VERSION=testing \
CLICKHOUSE_ENABLED=false \
USE_TLS=false \
USE_JWT=false \
NOMINATIM_ADDRESS=https://nominatim.openstreetmap.org/search \
LOG_LEVEL=debug \
WEBSERVER_LISTEN_ADDRESS=:8001 \
go run ./cmd/main.go
```

This started the API at `http://127.0.0.1:8001`.

---

## 3. Slovenian Normalization Implementation

We created a new normalizer for Slovenia (`internal/normalization/si.go`) and corresponding tests (`internal/normalization/si_test.go`).

- **Normalization handles:**
  - Slovenian postal code formats (4 digits)
  - Slovenian diacritics
  - Common street abbreviations (e.g., `ul.` → `ulica`, `c.` → `cesta`)
  - Special character removal

Example test cases:
```go
{name: "Valid 4-digit postal code", input: "1000", expect: "1000"},
{name: "Street with abbreviation", input: "Slovenska ul.", expect: []string{"slovenska ulica"}},
{name: "City with diacritics", input: "Črnomelj", expect: "crnomelj"},
```

---

## 4. Running and Fixing the Tests

We ran the Slovenian normalization tests with:

```bash
go test ./internal/normalization -v -run SI
```

### Issues and Fixes
- **Postal code normalization:**
  - Fixed extraction of the first 4 digits, ignoring extra characters
  - Adjusted error handling for too-short codes
- **City normalization:**
  - Improved special character removal and trimming
- **All tests now pass** for Slovenian normalization.

---

## 5. API Testing with Real Slovenian Addresses

We tested the running OpenPAQ API with real Slovenian addresses using `curl`:

### Example: Valid Address (Ljubljana)
```bash
curl "http://127.0.0.1:8001/api/v1/check?street=Slovenska+cesta+1&postal_code=1000&city=Ljubljana&country_code=SI&debug_details=true"
```
**Result:**
```json
{
  "street": "Slovenska cesta 1",
  "city": "Ljubljana",
  "postal_code": "1000",
  "country_code": "si",
  "street_matched": true,
  "city_matched": true,
  "postal_code_matched": true,
  "city_to_postal_code_matched": true,
  "country_code_matched": true,
  "version": "testing"
}
```

### Example: Valid Address (Maribor)
```bash
curl "http://127.0.0.1:8001/api/v1/check?street=Glavni+trg+8&postal_code=2000&city=Maribor&country_code=SI&debug_details=false"
```
**Result:**
All fields matched as expected.

### Example: Invalid Address
```bash
curl "http://127.0.0.1:8001/api/v1/check?street=Invalid+Street+999&postal_code=9999&city=InvalidCity&country_code=SI&debug_details=false"
```
**Result:**
All match fields returned `false`.

---

## 6. Interactive Bash Script for Slovenian Address Testing

To make testing easier, we created a simple interactive bash script: `test_slovenian_address.sh`.

### Script Features
- Accepts a Slovenian address in natural format (e.g., `Tehnološki park 21, 1000 Ljubljana`)
- Parses the address into street, postal code, and city
- Formats and sends the API request to OpenPAQ
- Displays a color-coded summary of the validation results
- Prints the raw API response for further inspection

### Usage
Make the script executable:
```bash
chmod +x test_slovenian_address.sh
```

Run the script with a Slovenian address:
```bash
./test_slovenian_address.sh "Tehnološki park 21, 1000 Ljubljana"
```

### Example Output
```
=== OpenPAQ Slovenian Address Validator ===
Input address: Tehnološki park 21, 1000 Ljubljana

Parsed components:
  Street: Tehnološki park 21
  Postal Code: 1000
  City: Ljubljana
  Country: SI

Making API call to OpenPAQ...

=== Validation Results ===
  Street: ✓ Valid
  City: ✓ Valid
  Postal Code: ✓ Valid
  City-Postal Match: ✓ Valid
  Country: ✓ Valid

=== Raw API Response ===
{
    "street": "Tehnolo\u0161ki park 21",
    "city": "Ljubljana",
    "postal_code": "1000",
    "country_code": "si",
    "street_matched": true,
    "city_matched": true,
    "postal_code_matched": true,
    "city_to_postal_code_matched": true,
    "country_code_matched": true,
    "version": "testing"
}
```

You can use this script for any Slovenian address in the format:
```
Street [Number], PostalCode City
```
For example:
- `Ulica 7. Maja 4, 6250 Ilirska Bistrica`
- `Slovenska cesta 1, 1000 Ljubljana`
- `Glavni trg 8, 2000 Maribor`

---

## 6. Summary of Findings

- **Slovenian normalization** is now implemented and tested.
- **All unit tests pass** for Slovenian address normalization.
- **API validation** works for real Slovenian addresses, matching valid addresses and rejecting invalid ones.
- **Reference database:** All validation is performed against OpenStreetMap data via the Nominatim API.

---

## 7. Next Steps
- For a full-scale benchmark, create a dataset of 90 valid and 40 invalid Slovenian addresses and run batch validation.
- Compare results to ground truth for accuracy metrics (not yet implemented/tested). 