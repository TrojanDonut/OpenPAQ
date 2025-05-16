OpenPAQ (**Open P**ostal **A**ddress **Q**uality) is a tool to validate parts of a postal address. Checks are done via an HTTP-Endpoint.

The following data will be checked:

- street
- city
- postal code
- country code


OpenPAQ uses [Nominatim](https://github.com/osm-search/Nominatim) in the background to check addresses. OpenPAQ is tested against a self-hosted version from [mediagis](https://github.com/mediagis/nominatim-docker).

Please have a look at the [documentation](https://openpaq.de) for a detailed description of the program.
