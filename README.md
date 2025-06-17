OpenPAQ (**Open P**ostal **A**ddress **Q**uality) is a tool designed to validate key components of postal addresses, ensuring data accuracy and integrity. A core capability of OpenPAQ is the normalization of addresses as they are found "in the wild." It transforms diverse and inconsistent address inputs into a standardized format that is valid for performing checks against the OpenStreetMap database.

Operating via a simple HTTP-Endpoint, OpenPAQ leverages [Nominatim](https://github.com/osm-search/Nominatim) in the background to perform its checks. It is tested against a self-hosted Nominatim version from [mediagis](https://github.com/mediagis/nominatim-docker)..

OpenPAQ checks the following address components:


- street
- city
- postal code
- country code

### Key Features of OpenPAQ:

OpenPAQ offers the following capabilities in address validation, measured with an internal benchmark of 90 correct and 40 incorrect postal addresses per country:
1.	International Address Validation: Provides address validation with accuracy levels of approximately 80% or higher for key European countries (DE, NL, AT, CH, FR, GB, IT, PL, DK)
2.	Correct Address Identification: Achieves a recall rate of approximately 75% or higher for most benchmarked countries in identifying correct addresses.
3.	Incorrect Address Detection: Offers an F1 score of over 75% for most benchmarked countries in identifying incorrect addresses.


Please have a look at the [documentation](https://openpaq.de) for a detailed description of the program.
