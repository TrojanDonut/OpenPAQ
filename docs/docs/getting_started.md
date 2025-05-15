# Getting started

See how to start an OpenPAQ instance with a minimalistic configuration and all optional features disabled.
Nominatim (openstreetmap.org) is used as the backend for the address to be checked.

!!! warning
    * Please read the usage [policies](https://operations.osmfoundation.org/policies/nominatim/) of openstreetmap.org before starting. 
    * This configuration is only for demonstration purposes. 
    * Do not use in a production environment. 

### Prerequisites

Install [Docker](https://www.docker.com/) to run a container version of Golang if you want to compile and run the minimal configuration.

Install [Golang](https://go.dev/) when you want to compile and run yourself.

#### optional

- [Clickhouse](https://clickhouse.com/)
- [Memcached](https://memcached.org/)
- [OpenSSL](https://www.openssl.org/)



## Run the project

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


#### Minimal configuration

```shell
docker run \
-e CACHE_ENABLED=false \
-e VERSION=testing \
-e CLICKHOUSE_ENABLED=false \
-e USE_TLS=false \
-e USE_JWT=false \
-e NOMINATIM_ADDRESS=https://nominatim.openstreetmap.org/search \
-e LOG_LEVEL=debug \
-e WEBSERVER_LISTEN_ADDRESS=:8001 \
-p 8001:8001 \
ghcr.io/DENICeG/OpenPAQ:latest
```


#### Example Request

Run your request against the new OpenPAQ instance:

=== "without debug details"
    ```shell
    curl http://127.0.0.1:8001/api/v1/check\
    ?street=Theodor-Stern-Kai\
    &postal_code=60596\
    &city=Frankfurt\
    &country_code=DE\
    &debug_details=false
    ```

=== "with debug details"
    ```shell
    curl http://127.0.0.1:8001/api/v1/check\
    ?street=Theodor-Stern-Kai\
    &postal_code=60596\
    &city=Frankfurt\
    &country_code=DE\
    &debug_details=true
    ```


#### Example Response

Congratulation you should get your first response which should look similar to this:

=== "without debug details"

    ```json
    {
        "street": "Theodor-Stern-Kai",
        "city": "Frankfurt",
        "postal_code": "60596",
        "country_code": "de",
        "street_matched": true,
        "city_matched": true,
        "postal_code_matched": true,
        "city_to_postal_code_matched": true,
        "country_code_matched": true,
        "version": "testing",
        
    }
    ```

=== "with debug details"

    ```json
    {
        "street": "Theodor-Stern-Kai",
        "city": "Frankfurt",
        "postal_code": "60596",
        "country_code": "de",
        "street_matched": true,
        "city_matched": true,
        "postal_code_matched": true,
        "city_to_postal_code_matched": true,
        "country_code_matched": true,
        "version": "testing",
        "details": {
            "parameters": {
                "matching_algorithm": 0,
                "matching_threshold": 0,
                "AllowPartialMatch": false,
                "AllowPartialCompareListMatch": false,
                "partial_matching_algorithm": 0,
                "partial_matching_threshold": 0,
                "PartialInputSeparators": null,
                "PartialExcludeWords": null,
                "PartialCompareListSeparators": null,
                "AllowCombineAllForwardCombinations": false,
                "allowed_amount_of_changed_chars": 0
            },
            "city_street_matches": [
                {
                    "city": "frankfurt am main",
                    "street": "theodor stern kai",
                    "postal_code": "60596",
                    "country_code": "de",
                    "street_similarity": 1,
                    "was_partial_street_match": false,
                    "city_similarity": 0.5294118,
                    "was_partial_city_match": true,
                    "was_list_match": false
                },
                {
                    "city": "frankfurt am main",
                    "street": "theodor stern kai",
                    "postal_code": "60528",
                    "country_code": "de",
                    "street_similarity": 1,
                    "was_partial_street_match": false,
                    "city_similarity": 0.5294118,
                    "was_partial_city_match": true,
                    "was_list_match": false
                },
                {
                    "city": "frankfurt",
                    "street": "theodor stern kai",
                    "postal_code": "60596",
                    "country_code": "de",
                    "street_similarity": 1,
                    "was_partial_street_match": false,
                    "city_similarity": 1,
                    "was_partial_city_match": false,
                    "was_list_match": false
                },
                {
                    "city": "frankfurt",
                    "street": "theodor stern kai",
                    "postal_code": "60528",
                    "country_code": "de",
                    "street_similarity": 1,
                    "was_partial_street_match": false,
                    "city_similarity": 1,
                    "was_partial_city_match": false,
                    "was_list_match": false
                }
            ],
            "postal_code_street_matches": [
                {
                    "postal_code": "60596",
                    "street": "theodor stern kai",
                    "country_code": "de",
                    "street_similarity": 1,
                    "was_partial_street_match": false,
                    "was_list_match": false
                }
            ],
            "city_postal_code_matches": [
                {
                    "city": "frankfurt am main",
                    "postal_code": "60596",
                    "country_code": "de",
                    "city_similarity": 0.5294118,
                    "was_partial_city_match": true,
                    "was_list_match": false
                },
                {
                    "city": "frankfurt",
                    "postal_code": "60596",
                    "country_code": "de",
                    "city_similarity": 1,
                    "was_partial_city_match": false,
                    "was_list_match": false
                }
            ]
        }
    }
    ```
