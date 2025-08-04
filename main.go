package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const (
	kmPerDegree = 111.0 // Approximate number of kilometers per degree of latitude
	maxRadiusKm = 25.0  // Maximum radius for search
	gridStepKm  = 2.5   // Distance between grid points in kilometers
	maxWorkers  = 4     // Maximum number of concurrent workers
)

type SearchParams struct {
	Latitude   float64
	Longitude  float64
	Query string
	RadiusKm   float64
}

type Place struct {
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Stars       float64     `json:"rating"`
	Reviews     int         `json:"reviews"`
	Coordinates Coordinates `json:"location"`
	Hours       string      `json:"hours,omitempty"`
	Phone       string      `json:"phone,omitempty"`
	Website     string      `json:"website,omitempty"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}


var (
	latitude   float64
	longitude  float64
	searchTerm string
	radiusKm   float64
)

var rootCmd = &cobra.Command{
	Use:   "mapsscrap",
	Short: "A Google Maps business scraper",
	Long: `mapsscrap is a CLI tool that scrapes business information from Google Maps 
using web automation. It collects details like business names, addresses, 
ratings, review counts, and phone numbers for a given search term and location.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := SearchParams{
			Latitude:   latitude,
			Longitude:  longitude,
			Query: searchTerm,
			RadiusKm:   radiusKm,
		}

		if params.RadiusKm > maxRadiusKm {
			return fmt.Errorf("radius too large, maximum is %v km", maxRadiusKm)
		}

		runSearch(params)
		return nil
	},
}

func init() {
	rootCmd.Flags().Float64VarP(&latitude, "lat", "a", 0, "Latitude of search center")
	rootCmd.Flags().Float64VarP(&longitude, "lon", "o", 0, "Longitude of search center")
	rootCmd.Flags().StringVarP(&searchTerm, "query", "t", "", "Search query")
	rootCmd.Flags().Float64VarP(&radiusKm, "radius", "r", 2.0, "Search radius in kilometers")

	rootCmd.MarkFlagRequired("lat")
	rootCmd.MarkFlagRequired("lon")
	rootCmd.MarkFlagRequired("query")
}

func main() {
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runSearch(params SearchParams) {
	if params.RadiusKm > maxRadiusKm {
		fmt.Println("Radius too large, maximum is 25km")
		return
	}

	gridPoints := generateSearchGrid(
		params.Latitude,
		params.Longitude,
		params.RadiusKm,
		gridStepKm,
	)
	text := fmt.Sprintf("Searching %d locations in a radius of %.1f km around (%.6f, %.6f) for query '%s'.",
		len(gridPoints), params.RadiusKm, params.Latitude, params.Longitude, params.Query)
	fmt.Println(text)

	estimatedTime := estimateJobTime(len(gridPoints), maxWorkers)
	barText := fmt.Sprintf("Please wait... Estimated time: %s", estimatedTime)
	bar := progressbar.Default(int64(len(gridPoints)), barText)

	maxWorkers := maxWorkers
	results := make(chan []Place, len(gridPoints))
	var wg sync.WaitGroup

	// Process grid points in batches
	for i := 0; i < len(gridPoints); i += maxWorkers {
		end := i + maxWorkers
		if end > len(gridPoints) {
			end = len(gridPoints)
		}

		// Launch workers for this batch
		for j := i; j < end; j++ {
			wg.Add(1)
			params := SearchParams{
				Latitude:   gridPoints[j].Lat,
				Longitude:  gridPoints[j].Lon,
				Query: params.Query,
				RadiusKm:   1.0,
			}

			go searchWorker(params, results, &wg, bar)
		}

		// Wait for batch to complete
		wg.Wait()
		time.Sleep(2 * time.Second) // Rate limiting between batches
	}

	// Collect all results
	allPlaces := make([]Place, 0)
	close(results)

	// Process results and remove duplicates
	for places := range results {
		for _, place := range places {
			if !containsPlace(allPlaces, place) {
				allPlaces = append(allPlaces, place)
			}
		}
	}

	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error getting working directory: %v", err))
	}

	// Create data directory if it doesn't exist
	// Create save path with timestamp
	savePath := filepath.Join(workDir,
		fmt.Sprintf("prospects_%s.csv", time.Now().Format("2006-01-02_15-04-05")))

	if err := savePlacesToCSV(allPlaces, savePath); err != nil {
		panic(fmt.Sprintf("Error saving places to CSV: %v", err))
	}
	fmt.Printf("%d places saved to %s\n", len(allPlaces), savePath)
}

func estimateJobTime(numTasks int, maxWorkers int) time.Duration {
    if numTasks <= 0 {
        return 0
    }

    const taskDuration = 1 * time.Minute

    // If tasks are less than or equal to max workers, only one batch needed
    if numTasks <= maxWorkers {
        return taskDuration
    }

    // Calculate number of batches needed
    numBatches := int(math.Ceil(float64(numTasks) / float64(maxWorkers)))
    
    // Each batch takes taskDuration + 2 seconds pause between batches
    totalTime := time.Duration(numBatches) * (taskDuration + 2*time.Second)

    return totalTime
}

func searchWorker(params SearchParams, results chan<- []Place, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	defer bar.Add(1)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// Create done channel for timeout handling
	done := make(chan bool)
	var places []Place
	var err error

	// Run scraping in goroutine
	go func() {
		places, err = scrapeGoogleMaps(params)
		if err != nil {
			fmt.Printf("Error searching at point %.6f, %.6f: %v\n", params.Latitude, params.Longitude, err)
		}
		done <- true
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		if err == nil {
			results <- places
		}
	case <-ctx.Done():
		fmt.Printf("Search timed out for coordinates: %.6f, %.6f\n", params.Latitude, params.Longitude)
	}
}

func containsPlace(places []Place, newPlace Place) bool {
	for _, p := range places {
		if p.Name == newPlace.Name && p.Address == newPlace.Address {
			return true
		}
	}
	return false
}

func generateSearchGrid(centerLat, centerLng float64, radiusKm float64, stepKm float64) []Coordinates {
	// Calculate degree deltas
	latDelta := radiusKm / kmPerDegree
	// Longitude degrees per km varies with latitude
	lngDelta := radiusKm / (kmPerDegree * math.Cos(centerLat*math.Pi/180.0))

	// Calculate steps
	latSteps := int(math.Ceil(2 * radiusKm / stepKm))
	lngSteps := int(math.Ceil(2 * radiusKm / stepKm))

	// Generate grid points
	points := make([]Coordinates, 0, latSteps*lngSteps)

	for i := 0; i < latSteps; i++ {
		for j := 0; j < lngSteps; j++ {
			lat := centerLat - latDelta + (2 * latDelta * float64(i) / float64(latSteps-1))
			lon := centerLng - lngDelta + (2 * lngDelta * float64(j) / float64(lngSteps-1))
			points = append(points, Coordinates{Lat: lat, Lon: lon})
		}
	}

	return points
}

func scrapeGoogleMaps(params SearchParams) ([]Place, error) {
	// Launch browser
	launch := launcher.New().
		Headless(true).
		Devtools(false)

	url, err := launch.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.Close()

	page := browser.MustPage()
	defer page.Close()

	// Navigate to Google Maps
	mapURL := fmt.Sprintf("https://www.google.com/maps/search/%s/@%f,%f,15z",
		params.Query,
		params.Latitude,
		params.Longitude,
	)

	if err := page.Navigate(mapURL); err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	page.MustWaitStable()

	listDivClass := "m6QErb.DxyBCb.kA9KIf.dS8AEf"
	places := []Place{}

	container := page.MustElement("div." + listDivClass)
	container.MustWaitVisible()

	// move mouse pointer to list which is first third of screen and scroll
	for i := 0; i < 10; i++ { // 10
		page.Mouse.MoveTo(proto.Point{X: 250, Y: 300})
		page.Mouse.Scroll(0.0, 6000.0, 30)
		// page.Mouse.Scroll(0.0, 1000.0, 5)
		time.Sleep(500 * time.Millisecond)
	}

	placeElements := container.MustElements("div.Nv2PK")

	for _, element := range placeElements {
		place := extractPlaceDetails(element, params)
		if place.Name != "" {
			places = append(places, place)
		}
	}

	return places, nil
}

func extractPlaceDetails(element *rod.Element, params SearchParams) Place {
	place := Place{
		Coordinates: Coordinates{
			Lat: params.Latitude,
			Lon: params.Longitude,
		},
	}

	// Extract place details
	if nameEl, err := element.Element("div.qBF1Pd.fontHeadlineSmall"); err == nil {
		place.Name = nameEl.MustText()
	}

	if ratingEl, err := element.Element("span.MW4etd"); err == nil {
		ratingText := ratingEl.MustText()
		fmt.Sscanf(ratingText, "%f", &place.Stars)
	}

	if reviewsEl, err := element.Element("span.UY7F9"); err == nil {
		reviewText := reviewsEl.MustText()
		fmt.Sscanf(reviewText, "(%d)", &place.Reviews)
	}

	if addressEl, err := element.Element("div.W4Efsd:nth-child(1)"); err == nil {
		line, err := addressEl.Text()
		if err == nil {
			lineSplit := strings.Split(line, "·")
			address := lineSplit[len(lineSplit)-1]
			place.Address = address
		}
	}

	if oppeningHoursEl, err := element.Element("div.W4Efsd:nth-child(2)"); err == nil {
		line, err := oppeningHoursEl.Text()
		if err == nil {
			lineSplit := strings.Split(line, "·")
			if len(lineSplit) > 1 {
				openingHours := lineSplit[0]
				place.Hours = openingHours
			}
		}
	}

	if phoneEl, err := element.Element("div.W4Efsd span.UsdlK"); err == nil {
		phone, err := phoneEl.Text()
		if err == nil {
			place.Phone = phone
		}
	}

	if websiteEl, err := element.Element("a.lcr4fd"); err == nil {
		if href, err := websiteEl.Attribute("href"); err == nil {
			place.Website = *href
		}
	}

	return place
}

func savePlacesToCSV(places []Place, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Name", "Address", "Stars", "Reviews", "Phone", "Hours", "Website"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header to CSV: %w", err)
	}

	// Write place data
	for _, place := range places {
		record := []string{
			place.Name,
			place.Address,
			fmt.Sprintf("%.1f", place.Stars),
			fmt.Sprintf("%d", place.Reviews),
			place.Phone,
			place.Hours,
			place.Website,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to CSV: %w", err)
		}
	}

	return nil
}
