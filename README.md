# Google Maps Scraper

Mapsscraper is a simple and free open-source CLI tool for scraping business information from Google Maps.

Built in less than 500 lines of Go. The application scrapes business information from Google Maps using web automation. Given a set of coordinates, a query string and a search raduis, the tool collects details like business names, addresses, ratings, review counts, and phone numbers for a given search term and location and saves the results in a CSV file. A tipical search will find a couple of hundred results in a few minutes.

## Features

- 🌎 Search businesses by latitude/longitude coordinates
- 📍 Configurable search radius
- 🤖 Automated browser interaction
- 💾 CSV export
- 📱 Collects Business Details: Business Name, Address, Stars, # of reviews, Phone Number and Website

## Example Result Table

An example result can found in the 'example_results' directory.

To search for prosthodontists in Mexico City within a 5 km radius, you can use the following command:
```
mapsscrap --lat 19.4343491 --lon -99.1775742 --query "prosthodontists" --radius 20
```

| Name | Address | Stars | Reviews | Phone | Hours | Website |
| --- | --- | --- | --- | --- | --- | --- |
| Dental Services Lomas | Perif. Blvd. Manuel Ávila Camacho 66-Primer Piso | 5 | 3 | 55 4952 2356 | Abierto ⋅ Cierra a las 7 p.m. | dentalslomas.com.mx |
| Alicama 16 Di Dental Lomas | Alicama 16-1 piso, int 5 | 5 | 3 | 55 6905 6286 | Cerrado ⋅ Abre a las 12 p.m. | www.doctoralia.com.mx|
| Dental Moliere | Av. Emilio Castelar 240 | 5 | 9 | 55 4132 5981 | Abierto ⋅ Cierra a las 7 p.m. | www.dentalmoliere.com.mx |


## Installation

Releases can be found on the [GitHub Releases page](https://github.com/edlgg/mapsscrap/releases).

If you are not sure how this works copy the content of this file and of the release page to you LLM of choice and ask it how to do it.
It should take 2 or 3 commands at the terminal to get it running.

### Linux Example
```bash
# One-line installation (x86_64)
curl -L https://github.com/edlgg/mapsscrap/releases/latest/download/mapsscrap_linux_x86_64 -o mapsscrap && chmod +x mapsscrap && sudo mv mapsscrap /usr/local/bin/

# Or for ARM64
curl -L https://github.com/edlgg/mapsscrap/releases/latest/download/mapsscrap_linux_arm64 -o mapsscrap && chmod +x mapsscrap && sudo mv mapsscrap /usr/local/bin/
```

### Usage
```bash
mapsscrap --help
mapsscrap --lat 19.4343491 --lon -99.1775742 --query "lawyer" --radius 20
```

## Roadmap
- [x] Scrape each element in list
- [x] Save results to CSV
- [x] Add scrolling to load more results
- [x] parse all wanted elements
- [x] get website for each place
- [x] add radius search in klm
- [x] parallelize
- [x] headless
- [x] add progress bar
- [x] convert to cli
- [x] research way to easy install
- [ ] add coordinates and query to results
- [ ] env var so only one execution can run at a time
- [ ] add logging for errors and warnings
- [ ] add example query and result
- [ ] make readme more user focus. "Want to get hundreds of potential leads in a couple of minutes? Use this tool to scrape Google Maps for business information."
- [ ] add comments to methods
- [ ] add save_path argument. Default should be ./out.csv
- [ ] add integration tests
- [ ] add unit tests
- [ ] functionality for multiple queries and locations