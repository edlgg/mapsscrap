# Google Maps Scraper

Mapsscraper is a simple and free open-source CLI tool for scraping business information from Google Maps.

Built in less than 500 lines of Go. The application scrapes business information from Google Maps using web automation. Given a set of coordinates, a query string and a search raduis, the tool collects details like business names, addresses, ratings, review counts, and phone numbers for a given search term and location and saves the results in a CSV file.

## Features

- 🌎 Search businesses by latitude/longitude coordinates
- 📍 Configurable search radius
- 🤖 Automated browser interaction
- 💾 CSV export
- 📱 Collects business details:
  - Business name
  - Address
  - Stars
  - Number of reviews
  - Phone number
  - Website

## Installation

Releases can be found on the [GitHub Releases page](https://github.com/edlgg/mapsscrap/releases).

### Linux
```bash
# One-line installation (x86_64)
curl -L https://github.com/edlgg/mapsscrap/releases/latest/download/mapsscrap_linux_x86_64 -o mapsscrap && chmod +x mapsscrap && sudo mv mapsscrap /usr/local/bin/

# Or for ARM64
curl -L https://github.com/edlgg/mapsscrap/releases/latest/download/mapsscrap_linux_arm64 -o mapsscrap && chmod +x mapsscrap && sudo mv mapsscrap /usr/local/bin/
```

### Usage
```bash
mapsscrap --help
mapsscrap --lat 19.4343491 --lon -99.1775742 --query "lawyer" --radius 5
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
- [ ] add logging for errors and warnings
- [ ] add save_path argument. Default should be ./out.csv
- [ ] add integration tests
- [ ] add unit tests