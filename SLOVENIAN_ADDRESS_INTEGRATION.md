# Slovenian Address Integration

## Overview

This document explains the modifications made to the OpenPAQ project to replace Nominatim/OpenStreetMap address checking with a local SQLite database containing all Slovenian addresses. This change provides faster, more reliable address validation for Slovenian addresses without external API dependencies.

## Changes Made

### 1. CSV to SQLite Conversion

**File**: `scripts/convert_csv_to_sqlite.go`

Created a conversion script that:
- Reads the Slovenian addresses CSV file (`RN_SLO_NASLOVI_register_naslovov_20250817.csv`)
- Creates a SQLite database (`slovenian_addresses.db`) with optimized schema
- Imports ~1.09 million address records
- Creates indexes for efficient querying on:
  - `obcina_naziv` (municipality name)
  - `naselje_naziv` (settlement/city name) 
  - `ulica_naziv` (street name)
  - `postni_okolis_naziv` (postal area name)
  - `postni_okolis_sifra` (postal code)

**Database Schema**:
```sql
CREATE TABLE slovenian_addresses (
    feature_id TEXT PRIMARY KEY,
    eid_naslov TEXT,
    obcina_sifra TEXT,
    obcina_naziv TEXT,
    naselje_sifra TEXT,
    naselje_naziv TEXT,
    ulica_sifra TEXT,
    ulica_naziv TEXT,
    postni_okolis_sifra TEXT,
    postni_okolis_naziv TEXT,
    hs_stevilka TEXT,
    hs_dodatek TEXT,
    -- ... additional fields
);
```

### 2. New SQLite Address Checker

**File**: `internal/slodb/db.go`

Created a new package that implements the same interface as the Nominatim checker but uses SQLite:

**Key Components**:
- `SIAddressDB` struct: Manages SQLite connection and configuration
- `Handle()` method: Main entry point that orchestrates address checking
- `CityStreetCheck()`: Validates city and street combinations
- `PostalCodeStreetCheck()`: Validates postal code and street combinations  
- `PostalCodeCityCheck()`: Validates postal code and city combinations

**Query Strategy**:
1. **Exact Match**: First tries exact case-insensitive matches
2. **Fuzzy Match**: Falls back to LIKE-based queries with fuzzy string matching
3. **Similarity Scoring**: Uses the same `algorithms.GetMatches()` as the original system

### 3. Service Integration

**Files Modified**:
- `internal/service.go`: Added `SIAddressDB` field and initialization
- `internal/httpserver.go`: Modified request handling to use SI DB for `country_code=si`
- `cmd/main.go`: Added `SI_ADDRESSES_DB_PATH` environment variable support

**Request Flow**:
```
Input Request → Check country_code
├── country_code = "si" → Use SQLite DB (SIAddressDB)
└── country_code ≠ "si" → Use Nominatim (existing behavior)
```

### 4. Environment Configuration

**New Environment Variable**:
- `SI_ADDRESSES_DB_PATH`: Path to the SQLite database file (optional)

**Example Usage**:
```bash
export SI_ADDRESSES_DB_PATH=/path/to/slovenian_addresses.db
```

## Database Statistics

After conversion, the SQLite database contains:
- **Total Records**: 1,093,570 addresses
- **Unique Cities**: 5,297 settlements
- **Unique Streets**: 6,279 street names  
- **Unique Postal Codes**: 466 postal codes

## Performance Benefits

### Before (Nominatim)
- External API calls with network latency
- Rate limiting and potential service outages
- Dependency on internet connectivity
- Slower response times (typically 200-500ms)

### After (SQLite)
- Local database queries (typically <50ms)
- No external dependencies
- Consistent availability
- Better scalability for high-volume requests

## API Compatibility

The API interface remains unchanged. All existing endpoints and response formats are preserved:

**Request**:
```
GET /api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si
```

**Response**:
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

## Setup Instructions

### 1. Convert CSV to SQLite
```bash
# Run the conversion script
go run scripts/convert_csv_to_sqlite.go
```

### 2. Set Environment Variable
```bash
export SI_ADDRESSES_DB_PATH=$(pwd)/slovenian_addresses.db
```

### 3. Start the Server
```bash
go run ./cmd
```

### 4. Test Slovenian Address
```bash
curl "http://localhost:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si"
```

## Data Source

The Slovenian address data comes from the official Register of Addresses (Register naslovov) maintained by the Surveying and Mapping Authority of the Republic of Slovenia (GURS). The CSV file contains the complete, up-to-date database of all Slovenian addresses as of the export date.

## Dependencies Added

- `github.com/mattn/go-sqlite3`: SQLite driver for Go

## Future Enhancements

Potential improvements for the Slovenian address checker:
1. **Caching**: Add in-memory caching for frequently queried addresses
2. **Full-text Search**: Implement FTS5 for better fuzzy matching
3. **Geocoding**: Add coordinate-based queries using the E/N fields
4. **Updates**: Automated process to update the database from new CSV exports
5. **Compression**: Consider database compression for storage optimization

## Migration Notes

- **Backward Compatibility**: All existing functionality for non-Slovenian addresses remains unchanged
- **Gradual Rollout**: The change only affects `country_code=si` requests
- **Fallback**: If `SI_ADDRESSES_DB_PATH` is not set, Slovenian requests fall back to Nominatim
- **Monitoring**: Existing metrics and logging continue to work for the new SQLite path
