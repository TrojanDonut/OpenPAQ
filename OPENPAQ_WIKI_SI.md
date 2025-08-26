# OpenPAQ – Validacija naslovov (SI integracija)

## Kratek opis projekta
OpenPAQ je HTTP storitev za preverjanje, ali podani poštni naslovi obstajajo. Vhodni podatki (ulica, mesto, poštna številka, ISO koda države) se normalizirajo in nato preverijo preko različnih virov resnice. Rezultat je JSON s polji, ki povedo, ali so ulica, mesto, poštna številka, vezava mesto↔pošta in koda države ujemajoči.

## Kaj smo spremenili
- Dodali smo lokalno validacijo za slovenske naslove z uporabo SQLite baze namesto Nominatim/OpenStreetMap.
- Za `country_code=si` se zdaj uporablja lokalna baza z vsemi slovenskimi naslovi (hitro, brez omrežnih odvisnosti).
- Za ostale države je vedenje nespremenjeno in se uporablja Nominatim.

### Povzetek sprememb v kodi
- Dodana pretvorba CSV → SQLite:
  - `scripts/convert_csv_to_sqlite.go` (uvoz ~1.093.570 zapisov iz `RN_SLO_NASLOVI_*.csv` v `slovenian_addresses.db` + indeksi)
- Nov SI preverjalnik (SQLite):
  - `internal/slodb/db.go` z metodami `Handle`, `CityStreetCheck`, `PostalCodeStreetCheck`, `PostalCodeCityCheck`
- Integracija v storitev:
  - `internal/service.go`: nova lastnost `siDB`, konfiguracija `SIAddressesDBPath`
  - `internal/httpserver.go`: za `country_code=si` uporabi `siDB`, drugače Nominatim
  - `cmd/main.go`: dodan ENV `SI_ADDRESSES_DB_PATH`
- Dokumentacija in testiranje:
  - `SLOVENIAN_ADDRESS_INTEGRATION.md` (tehnična razlaga SI integracije)
  - `TESTING_GUIDE.md` (navodila za testiranje)
  - `test_slovenian_address.sh` (skripta za testni nabor – SI in globalni primeri)

## Arhitektura / potek
1. Odjemalec pokliče `GET /api/v1/check?...` z `street`, `city`, `postal_code`, `country_code`.
2. Vhod se normalizira (lowercase, čiščenje po državni logiki).
3. Odločitev vira resnice:
   - `country_code = si` in `SI_ADDRESSES_DB_PATH` je nastavljen → uporabi SQLite (`internal/slodb`).
   - sicer → Nominatim (obstoječe vedenje).
4. Izvedejo se tri neodvisne preverbe in sestavi rezultat:
   - Ujemanje Ulica–Mesto
   - Ujemanje Pošta–Ulica
   - Ujemanje Pošta–Mesto
5. Odziv vrne boolean polja in opcijsko “debug” detajle (če `debug_details=true`).

## Baza podatkov (Slovenija)
- Datoteka: `slovenian_addresses.db` (SQLite)
- Tabela: `slovenian_addresses`
- Ključni stolpci za validacijo: `naselje_naziv` (mesto/naselje), `ulica_naziv` (ulica), `postni_okolis_sifra` (poštna številka)
- Indeksi: na `naselje_naziv`, `ulica_naziv`, `postni_okolis_sifra` (hitrejše poizvedbe)
- Statistika (trenutni uvoz):
  - 1.093.570 zapisov, 5.297 unikatnih naselij, 6.279 ulic, 466 pošt

## Namestitev in zagon
1. Pretvori CSV v SQLite (enkratno ali po novi objavi CSV):
   ```bash
   go run scripts/convert_csv_to_sqlite.go
   ```
2. Zaženi strežnik z vključeno SI bazo:
   ```bash
   export SI_ADDRESSES_DB_PATH=$(pwd)/slovenian_addresses.db
   export LOG_LEVEL=debug
   export USE_TLS=false
   export USE_JWT=false
   export WEBSERVER_LISTEN_ADDRESS=127.0.0.1:8080
   export NOMINATIM_ADDRESS=https://nominatim.openstreetmap.org
   export VERSION=dev
   export CACHE_ENABLED=false
   export CLICKHOUSE_ENABLED=false
   go run ./cmd
   ```

## Testiranje
### Ročno (cURL)
- Slovenija (hitro, lokalna baza):
  ```bash
  curl "http://127.0.0.1:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si"
  ```
- Globalno (Nominatim):
  ```bash
  curl "http://127.0.0.1:8080/api/v1/check?street=Unter%20den%20Linden&city=Berlin&postal_code=10117&country_code=de"
  ```
- Debug način:
  ```bash
  curl "http://127.0.0.1:8080/api/v1/check?street=Trzinska%20cesta&city=Menge%C5%A1&postal_code=1234&country_code=si&debug_details=true"
  ```

### Avtomatizirano (skripta)
```bash
chmod +x test_slovenian_address.sh
./test_slovenian_address.sh
```
Skripta pokrije:
- ✅ Veljavne/neveljavne slovenske primere
- ✅ Veljavne/neveljavne globalne primere
- ✅ Robne primere in posebne znake
- ✅ Debug način
- ✅ Primerjavo odzivnih časov (SI <50 ms vs globalno 200–500 ms)

### Programatično (Go primeri)
```bash
go run test_examples.go
```

## Konfiguracija (ENV)
- `SI_ADDRESSES_DB_PATH` (opcijsko): pot do `slovenian_addresses.db`. Če ni nastavljeno, se tudi za SI uporabi Nominatim.
- Ostali (obstoječi): `WEBSERVER_LISTEN_ADDRESS`, `USE_TLS`, `TLS_CERT_FILE_PATH`, `TLS_KEY_FILE_PATH`, `USE_JWT`, `JWT_SIGNING_KEY`, `LOG_LEVEL`, `VERSION`, `CACHE_ENABLED`, `CACHE_URL`, `CLICKHOUSE_ENABLED`, ...

## Zmogljivost in zanesljivost
- **SI (SQLite)**: lokalno, brez omrežnih odvisnosti, tipičen odziv <50 ms.
- **Globalno (Nominatim)**: omrežje, omejitve in variabilni časi (200–500 ms).
- API ostane nespremenjen (isti endpointi in JSON polja).

## Vzdrževanje
- Ob novi izdaji CSV (Register naslovov – GURS):
  - Zaženi pretvorbo (`go run scripts/convert_csv_to_sqlite.go`) za posodobitev `slovenian_addresses.db`.
- Možne nadgradnje:
  - FTS5 (polno besedilno iskanje) za še boljšo približno ujemanje
  - Predpomnjenje vročih naslovov
  - Uporaba koordinat (stolpca E/N) za geografske preverbe
  - Avtomatizacija periodičnih posodobitev baze

## Primer odziva (SI)
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

---
Za več tehničnih detajlov glej `SLOVENIAN_ADDRESS_INTEGRATION.md` in `TESTING_GUIDE.md` v repozitoriju.
