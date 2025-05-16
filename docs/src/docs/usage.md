## Make a first request

```shell
    curl http://127.0.0.1:8001/api/v1/check\
    ?street=Köpenicker Landstr 252\
    &postal_code=12347\
    &city=Berlin Treptow\
    &country_code=DE\
    &debug_details=true
```

The response should look like this:
```json
{
  "street": "Max Mustermann,Koepenicker Landstr 252",
  "city": "Berlin Treptow",
  "postal_code": "12437",
  "country_code": "de",
  "street_matched": true,
  "city_matched": true,
  "postal_code_matched": true,
  "city_to_postal_code_matched": true,
  "country_code_matched": true,
  "version": "test-version"
}
```


## Interpretation of OpenPAQ's results

OpenPAQ is a powerful tool for verifying the existence of arbitrary addresses worldwide. During the validation process, input addresses are first preprocessed using a set of country-specific rules.

After preprocessing, each address is evaluated using three core contextual checks that help ensure the consistency and plausibility of its components:

## Validation Checks

### 1. Street–Postal Code Check

> Is there a street within the specified postal code area that approximately matches the given street name?

### 2. Street–City Check

> Does the specified city contain a street that approximately matches the given street name?

### 3. City–Postal Code Check

> Is there a city within the specified postal code area that approximately matches the given city name?

These checks leverage fuzzy matching to account for minor errors, variations, or misspellings in the input data.


---

## Interpretation of Validation Results

For detailed implementation, API usage, or customization options, please refer to the relevant sections in this documentation.

The results of the address validation checks can be interpreted as follows:

### `street_matched` (true / false)

Indicates whether the provided street is likely a valid component of the address. This value is derived by evaluating the following conditions:

- **Street–Postal Code Check**
- **Street–City Check**

If at least one of these conditions returns a match, the street is considered valid and `street_matched` is set to `true`.

---

### `postal_code_matched` (true / false)

Represents the likelihood that the given postal code is correct. This is determined using:

- **Street–Postal Code Check**
- **City–Postal Code Check**

If either condition holds true, the postal code is deemed valid and `postal_code_matched` is `true`.

---

### `city_matched` (true / false)

Indicates whether the provided city is a valid part of the address, based on the following checks:

- **Street–City Check**
- **City–Postal Code Check**

If at least one of the two conditions is fulfilled, `city_matched` will return `true`.

---

### `city_to_postal_code_matched` (true / false)

Even if both `postal_code_matched` and `city_matched` return `true`, this does not necessarily mean the city and postal code are directly related. Therefore, this field explicitly checks:

- **City–Postal Code Check**

This helps confirm that the city is indeed located within the provided postal code area.

---

### `country_code_matched` (true / false)

All validation checks are scoped to specific countries. This value returns `true` if any valid match found during the checks belongs to the same country code as the input address.

---

### Example: Street–City Check

The following example illustrates the logical evaluation of a Street–City Check using **Nominatim**.  
It demonstrates how preprocessing input data is crucial to handle additional elements frequently included in addresses—such as building names or extra descriptors.

The preprocessing phase helps normalize these variations, correct common typing errors, and isolate the core address components.  
This ensures the system can extract and validate the true underlying address information effectively.



``` mermaid

flowchart TD
    %%{init: {"flowchart": {"htmlLabels": false}} }%%

    A1[Input:
    - DE
    - Berlin Treptow
    - Max Mustermann, Koepenicker Landstr 252
    - 12437] --> A2[DE Normalizer]

    A2 --> A3[Normalized Input:
    - berlin treptow
    - max mustermann
    - koepenicker landstr
    - 12437]

    A3 --> A4.1[Nominatim Request:
    - search?q=berlin treptow, koepenicker landstr
    ⛳]

    A3 --> A4.2[Nominatim Request:
    - search?q=berlin treptow, max mustermann]


    A4.1 --> A5[Nominatim Results:
    - Köpenicker Landstraße
    - Berlin]

    A5 --> A6[Normalize Nominatim Results:
    - koepenicker landstrasse
    - berlin]

    A6 --> A7a[Fuzzy Comparison Street:
    Compare Input vs. Result
    **normalized input:** koepenicker landstr
    **normalized nominatim:** koepenicker landstrasse
    **Similarity:** 0.8636364]
    
    A6 --> A7b["Fuzzy Comparison City:
    Compare Input vs. Result
    **normalized input:** 
    berlin treptow
    **normalized nominatim:** berlin
    **Similarity:** 
    0.42857143 (partial match)"]

    A7a --> A8[Street-City-Check ✅]
    A7b --> A8[Street-City-Check ✅]

    A8 --> A9[Is street correct? ✅]
    A8 --> A10[Is city correct? ✅]

```


## Debug Information

#### parameter
Validates the setup parameters of the program, that are used in the address evaluation process. These parameters for the moment are in the code, they will be customizable at a later stage.

```json
{
"matching_algorithm": 1,
"matching_threshold": 0.8,
"AllowPartialMatch": false,
"AllowPartialCompareListMatch": false,
"partial_matching_algorithm": 1,
"partial_matching_threshold": 0.8,
"PartialInputSeparators": null,
"PartialExcludeWords": null,
"PartialCompareListSeparators": null,
"AllowCombineAllForwardCombinations": false,
"allowed_amount_of_changed_chars": 0
}
```
### city_street_matches
Contains all results for the Street-City-Check and important meta information.
The similarity values are the actual fuzzy similarities applying the utilized fuzzy distance algorithm between the result string and the initial input argument.
partial_match indicates if argument matches just matches for a separated part of the argument. For example input city "Frankfurt" would state
true for the result "Frankfurt am Main". For this reason, if partial matches are allowed for the given country, the similarity can be lower than the matching_threshold.
was_list_match is an indicator that states whether the result was extracted due to list matching or detected by nominatim. If no list matcher is configured, 
the field states always false and the result was found by nominatim.

#### Example output for `Max Mustermann, Koepenicker Landstr 252 12437 Berlin Treptow DE`
```json
[
  {
    "city": "berlin",
    "street": "koepenicker landstraße",
    "postal_code": "12437",
    "country_code": "de",
    "street_similarity": 0.8636364,
    "was_partial_street_match": false,
    "city_similarity": 0.42857143,
    "was_partial_city_match": true,
    "was_list_match": false
  },
  {
    "city": "treptow koepenick",
    "street": "koepenicker landstraße",
    "postal_code": "12437",
    "country_code": "de",
    "street_similarity": 0.8636364,
    "was_partial_street_match": false,
    "city_similarity": 0.23529412,
    "was_partial_city_match": true,
    "was_list_match": false
  },
  {
    "city": "berlin",
    "street": "koepenicker landstraße",
    "postal_code": "12435",
    "country_code": "de",
    "street_similarity": 0.8636364,
    "was_partial_street_match": false,
    "city_similarity": 0.42857143,
    "was_partial_city_match": true,
    "was_list_match": false
  },
  {
    "city": "treptow koepenick",
    "street": "koepenicker landstraße",
    "postal_code": "12435",
    "country_code": "de",
    "street_similarity": 0.8636364,
    "was_partial_street_match": false,
    "city_similarity": 0.23529412,
    "was_partial_city_match": true,
    "was_list_match": false
  },
  {
    "city": "berlin",
    "street": "koepenicker landstraße",
    "postal_code": "12439",
    "country_code": "de",
    "street_similarity": 0.8636364,
    "was_partial_street_match": false,
    "city_similarity": 0.42857143,
    "was_partial_city_match": true,
    "was_list_match": false
  },
  {
    "city": "treptow koepenick",
    "street": "koepenicker landstraße",
    "postal_code": "12439",
    "country_code": "de",
    "street_similarity": 0.8636364,
    "was_partial_street_match": false,
    "city_similarity": 0.23529412,
    "was_partial_city_match": true,
    "was_list_match": false
  }
]
```

### postal_code_street_matches
Contains all results for the Street-PostalCode-Check and important meta information.

#### Example output for `Max Mustermann, Koepenicker Landstr 252 12437 Berlin Treptow DE`
```json
[{
"postal_code": "12437",
"street": "koepenicker landstraße",
"country_code": "de",
"street_similarity": 0.8636364,
"was_partial_street_match": false,
"was_list_match": false
}]
```


### city_postal_code_matches
Contains all results for the City-PostalCode-Check and important meta information.

#### Example output for `Max Mustermann, Koepenicker Landstr 252 12437 Berlin Treptow DE`
```json
[{
"city": "berlin",
"postal_code": "12437",
"country_code": "de",
"city_similarity": 0.42857143,
"was_partial_city_match": true,
"was_list_match": false
},
{
"city": "treptow koepenick",
"postal_code": "12437",
"country_code": "de",
"city_similarity": 0.23529412,
"was_partial_city_match": true,
"was_list_match": false
}]
```