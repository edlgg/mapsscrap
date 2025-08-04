# Google Maps Scraper

A Go application that scrapes business information from Google Maps using web automation. Given a set of coordinates, a query string and a search raduis, the tool collects details like business names, addresses, ratings, review counts, and phone numbers for a given search term and location and saves the results in a CSV file.

## Features

- 🌎 Search businesses by latitude/longitude coordinates
- 📍 Configurable search radius
- 🤖 Automated browser interaction
- 💾 CSV export
- 📱 Collects business details:
  - Business name
  - Address
  - Rating
  - Number of reviews
  - Phone number

## Installation

### Linux/MacOS
```bash
# Download latest release (replace ARCH with amd64 or arm64)
curl -L https://github.com/edlgg/mapsscrap/releases/latest/download/mapsscrap_Linux_ARCH.tar.gz | tar xz

# Move to PATH
sudo mv mapsscrap /usr/local/bin/
```

### Windows
1. Download the latest release from https://github.com/edlgg/mapsscrap/releases
2. Extract the zip file
3. Add the executable to your PATH or run from the extracted location

### Usage
```bash
mapsscrap --help
```

## Roadmap
- [x] Scrape each element in list
- [x] Save results to CSV
- [x] Add scrolling to load more results
- [x] correct address
- [x] correct stars
- [x] get website for each place
- [x] get opening hours 
- [x] add radius search in klm
- [x] parallelize
- [x] headless
- [x] add progress bar
- [x] add logging
- [x] convert to cli
- [ ] research way to easy install
- [ ] add integration tests
- [ ] add unit tests