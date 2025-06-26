# Overview

OpenPAQ (**Open P**ostal **A**ddress **Q**uality) is a tool designed to validate key components of postal addresses, ensuring data accuracy and integrity. A core capability of OpenPAQ is the normalization of addresses as they are found "in the wild." It transforms diverse and inconsistent address inputs into a standardized format that is valid for performing checks against the OpenStreetMap database.

Operating via a simple HTTP-Endpoint, OpenPAQ leverages [Nominatim](https://github.com/osm-search/Nominatim) in the background to perform its checks. It is tested against a self-hosted Nominatim version from [mediagis](https://github.com/mediagis/nominatim-docker).



#### OpenPAQ checks the following address components:

- street
- city
- postal code
- country code

### Key Features of OpenPAQ:
OpenPAQ offers the following capabilities in address validation, measured with an internal benchmark of 90 correct and 40 incorrect postal addresses per country:

1. International Address Validation: Provides address validation with accuracy levels of approximately 80% or higher for key European countries (DE, NL, AT, CH, FR, GB, IT, PL, DK)
2. Correct Address Identification: Achieves a recall rate of approximately 75% or higher for most benchmarked countries in identifying correct addresses.
3. Incorrect Address Detection: Offers an F1 score of over 75% for most benchmarked countries in identifying incorrect addresses.

<br>
<div class="grid cards" markdown>

-   :material-clock-fast:{ .lg .middle } __Set up in 2 minutes__

    ---

    Download a container and set it up with a minimal instance

    [:octicons-arrow-right-24: Getting started](getting_started.md)

-   :material-cog-outline:{ .lg .middle } __Customize configuration__

    ---

    Adjust the configuration to your needs

    [:octicons-arrow-right-24: Configuration](configuration.md)

-   :material-chat-question-outline:{ .lg .middle } __Understanding OpenPAQ's responses__

    ---

    Get further information what the response of a request mean
    
    [:octicons-arrow-right-24: Interpretation of Results](usage.md)

-   :material-code-json:{ .lg .middle } __Internals__

    ---

    Here you get a deeper explanation about the code and how things are realized 

    [:octicons-arrow-right-24: Development](development.md)
  
-   :material-scale-balance:{ .lg .middle } __Contribute to the project__

    ---

    How to report bugs or make feature and pull requests

    [:octicons-arrow-right-24: Contribution](contributions.md)

- :material-scale-balance:{ .lg .middle } __License terms__

    ---

    OpenPAQ is licensed under AGPL
   

    [:octicons-arrow-right-24: License & EOL](license.md)

</div>






