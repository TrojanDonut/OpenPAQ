# Development


## build the project

```shell
go build  -o openPAQ -a ./cmd/main.go
```

## build the container
```shell
docker build -t openpaq:latest .
```

## run the project

```shell
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

## Program flow

``` mermaid
flowchart TD
    A1[User] --> A2[API Request]

    A2 --> A3["normlize input
    (country specific)"]

    A3 --> A4{ exists a country code specific listmatcher }

    subgraph LISTMATCHER_NOT_EXISTS
        S1_1[request nominatim] --> S1_2[normalize response]
        S1_2 --> S1_3[calculate match]
        S1_3 --> S1_4[response match]
        end
        
        subgraph LISTMATCHER_EXISTS
        
        S2_1[request nominatim] --> S2_2[normalize response]
        S2_2 --> S2_3[calculate match]
        
        S2_10[lookup in list] --> S2_11[normalize response]
        S2_11 --> S2_12[calculate match]
        
        S2_3 --> S20_1[response combined best match]
        S2_12 --> S20_1
        
    end
    
    A4 --- YES ---> LISTMATCHER_EXISTS
    A4 --- NO ---> LISTMATCHER_NOT_EXISTS
    
    LISTMATCHER_EXISTS-->A5[API Response]
    LISTMATCHER_NOT_EXISTS-->A5[API Response]
```

## Add a normalizer

### Create normalizer

A normalizer must implement following interface:

```go
type Normalize interface {
  GetCountryCode() string
  City(string) (string, error)
  PostalCode(string) (string, error)
  Street(string) ([]string, error)
}
```

A good starting point for a reference implementation is located at `internal/normalization/generic.go`.
This is a very minimal normalizer which will be used if no specific normalizer is registered for the provided 
country code.

### Register normalizer

To register a new normalizer, it needs do be added in the function `internal/normalization/normalize.go:NewNormalizer`

