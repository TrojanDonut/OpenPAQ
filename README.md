OpenPAQ (**OpenP**ostal **A**ddress **Q**uality) is a tool to validate parts of a postal address. Therefore, it provides an HTTP-API endpoint to provide the
relevant data.

Actually following data will be checked and validated:

- street
- city
- postal code
- country code


OpenPAQ use a Nominatim service in the background to check the addresses. Actually OpenPAQ is tested with a self-hosted
version and with the API from openstreetmap.org.

Please have a look in the [documentation](https://openpaq.de) for a detailed description of the program.