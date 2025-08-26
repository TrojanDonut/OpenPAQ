# OpenPAQ Testing Guide

This guide explains how to test both Slovenian addresses (using the SQLite database) and global addresses (using Nominatim) in the OpenPAQ system.

## Prerequisites

1. **Server Running**: Make sure the OpenPAQ server is running with Slovenian database support:
   ```bash
   export SI_ADDRESSES_DB_PATH=$(pwd)/slovenian_addresses.db
   go run ./cmd
   ```

2. **Database Ready**: Ensure the SQLite database has been created:
   ```bash
   go run scripts/convert_csv_to_sqlite.go
   ```

## Testing Methods

### 1. Manual Testing with cURL

#### Test Slovenian Addresses (SQLite Database)
```bash
# Valid Slovenian address
curl "http://127.0.0.1:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si"

# Another valid Slovenian address
curl "http://127.0.0.1:8080/api/v1/check?street=Cesta%20VIII&city=Gri%C4%8D&postal_code=1310&country_code=si"

# Invalid Slovenian address
curl "http://127.0.0.1:8080/api/v1/check?street=Nonexistent%20Street&city=Nonexistent%20City&postal_code=9999&country_code=si"
```

#### Test Global Addresses (Nominatim)
```bash
# Valid German address
curl "http://127.0.0.1:8080/api/v1/check?street=Unter%20den%20Linden&city=Berlin&postal_code=10117&country_code=de"

# Valid French address
curl "http://127.0.0.1:8080/api/v1/check?street=Champs-%C3%89lys%C3%A9es&city=Paris&postal_code=75008&country_code=fr"

# Valid UK address
curl "http://127.0.0.1:8080/api/v1/check?street=Oxford%20Street&city=London&postal_code=W1C%201AP&country_code=gb"

# Valid US address
curl "http://127.0.0.1:8080/api/v1/check?street=Times%20Square&city=New%20York&postal_code=10036&country_code=us"
```

#### Test with Debug Details
Add `&debug_details=true` to any request to see detailed matching information:
```bash
curl "http://127.0.0.1:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si&debug_details=true"
```

### 2. Automated Testing with Bash Script

Run the comprehensive test suite:
```bash
chmod +x test_slovenian_address.sh
./test_slovenian_address.sh
```

This script tests:
- ✅ Valid Slovenian addresses (SQLite)
- ✅ Invalid Slovenian addresses (SQLite)
- ✅ Valid global addresses (Nominatim)
- ✅ Invalid global addresses (Nominatim)
- ✅ Edge cases (empty fields, long names, special characters)
- ✅ Debug mode functionality
- ✅ Performance comparison

### 3. Programmatic Testing with Go

Run the Go test examples:
```bash
go run test_examples.go
```

This program demonstrates:
- Structured test cases for both address types
- JSON response parsing
- Performance measurement
- Error handling

## Expected Results

### Slovenian Addresses (SQLite Database)

**Valid Address Response**:
```json
{
  "street": "Trzinska cesta",
  "city": "Mengeš",
  "postal_code": "1234",
  "country_code": "si",
  "street_matched": true,
  "city_matched": true,
  "postal_code_matched": true,
  "city_to_postal_code_matched": true,
  "country_code_matched": true,
  "version": "dev"
}
```

**Performance**: Typically <50ms response time

### Global Addresses (Nominatim)

**Valid Address Response**:
```json
{
  "street": "Unter den Linden",
  "city": "Berlin",
  "postal_code": "10117",
  "country_code": "de",
  "street_matched": true,
  "city_matched": true,
  "postal_code_matched": true,
  "city_to_postal_code_matched": true,
  "country_code_matched": true,
  "version": "dev"
}
```

**Performance**: Typically 200-500ms response time

## Test Cases

### Slovenian Address Test Cases

| Test Case | Street | City | Postal Code | Expected Result |
|-----------|--------|------|-------------|-----------------|
| Valid 1 | Trzinska cesta | Mengeš | 1234 | ✅ Valid |
| Valid 2 | Cesta VIII | Grič | 1310 | ✅ Valid |
| Valid 3 | Šentviška pot | Čatež ob Savi | 8250 | ✅ Valid |
| Invalid | Nonexistent Street | Nonexistent City | 9999 | ❌ Invalid |

### Global Address Test Cases

| Test Case | Street | City | Postal Code | Country | Expected Result |
|-----------|--------|------|-------------|---------|-----------------|
| German | Unter den Linden | Berlin | 10117 | de | ✅ Valid |
| French | Champs-Élysées | Paris | 75008 | fr | ✅ Valid |
| UK | Oxford Street | London | W1C 1AP | gb | ✅ Valid |
| US | Times Square | New York | 10036 | us | ✅ Valid |
| Invalid | Fake Street 123 | Fake City | 00000 | us | ❌ Invalid |

## Performance Testing

### Response Time Comparison

```bash
# Test Slovenian address performance
time curl -s "http://127.0.0.1:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si" > /dev/null

# Test global address performance  
time curl -s "http://127.0.0.1:8080/api/v1/check?street=Unter%20den%20Linden&city=Berlin&postal_code=10117&country_code=de" > /dev/null
```

**Expected Results**:
- **Slovenian addresses**: <50ms (SQLite database)
- **Global addresses**: 200-500ms (Nominatim API)
- **Performance improvement**: ~5-10x faster for Slovenian addresses

## Debug Mode

Enable debug mode to see detailed matching information:

```bash
curl "http://127.0.0.1:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si&debug_details=true"
```

**Debug Response Includes**:
- Matching algorithm parameters
- Street-city matches with similarity scores
- Postal code matches
- Partial matching information

## Troubleshooting

### Common Issues

1. **Server not running**:
   ```bash
   # Check if server is running
   curl http://127.0.0.1:8080/version
   ```

2. **Database not found**:
   ```bash
   # Check if SQLite database exists
   ls -la slovenian_addresses.db
   ```

3. **Environment variable not set**:
   ```bash
   # Check environment variable
   echo $SI_ADDRESSES_DB_PATH
   ```

4. **Network issues (global addresses)**:
   ```bash
   # Test Nominatim connectivity
   curl https://nominatim.openstreetmap.org/
   ```

### Error Responses

**Database Connection Error**:
```json
{
  "error": "database connection failed"
}
```

**Invalid Parameters**:
```json
{
  "error": "street field exceed length limit of 500 elements"
}
```

**Service Unavailable**:
```json
{
  "error": "external service unavailable"
}
```

## Integration Testing

### Continuous Integration

For automated testing in CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Test OpenPAQ API
  run: |
    # Start server
    export SI_ADDRESSES_DB_PATH=./slovenian_addresses.db
    go run ./cmd &
    sleep 5
    
    # Run tests
    ./test_slovenian_address.sh
    
    # Run Go tests
    go run test_examples.go
```

### Load Testing

For performance testing:

```bash
# Install Apache Bench
sudo apt-get install apache2-utils

# Test Slovenian addresses
ab -n 100 -c 10 "http://127.0.0.1:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si"

# Test global addresses
ab -n 100 -c 10 "http://127.0.0.1:8080/api/v1/check?street=Unter%20den%20Linden&city=Berlin&postal_code=10117&country_code=de"
```

## Summary

- **Slovenian addresses**: Use SQLite database for fast, reliable validation
- **Global addresses**: Use Nominatim API for worldwide coverage
- **Same API interface**: Both address types use identical request/response format
- **Performance optimized**: Slovenian addresses are 5-10x faster
- **Comprehensive testing**: Multiple testing methods available
- **Debug support**: Detailed matching information available
