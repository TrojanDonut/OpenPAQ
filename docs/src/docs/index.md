# Overview

OpenPAQ (**OpenP**ostal **A**ddress **Q**uality) is a tool to validate parts of a postal address. Checks are done via an HTTP-Endpoint.

The following data will be checked:

- street
- city
- postal code
- country code


OpenPAQ uses [Nominatim](https://github.com/osm-search/Nominatim) in the background to check addresses. OpenPAQ is tested against a self-hosted version of the Nominatim docker container from [mediagis](https://github.com/mediagis/nominatim-docker).

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






