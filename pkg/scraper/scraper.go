package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/theoreotm/frieren-api/models"
	"github.com/theoreotm/frieren-api/pkg/data"
	"github.com/theoreotm/frieren-api/storage"
)

type Scraper struct {
	CharacterURLs        []string
	LocationURLs         []string
	URLSet               map[string]struct{}
	CharacterDataChannel chan *models.Character
	LocationDataChannel  chan *models.Location
	ShouldDebug          bool
	Logger               *logrus.Logger
}

type ScrapedData struct {
	CharacterURLs      []string
	AmountOfCharacters int
}

const (
	BaseFandomURL = "https://frieren.fandom.com"
)

// NewScraper initializes a new Scraper
func NewScraper(shouldDebug bool, logger *logrus.Logger) *Scraper {
	return &Scraper{
		CharacterURLs:        []string{},
		LocationURLs:         []string{},
		URLSet:               make(map[string]struct{}),
		CharacterDataChannel: make(chan *models.Character),
		LocationDataChannel:  make(chan *models.Location),
		ShouldDebug:          shouldDebug,
		Logger:               logger,
	}
}

func (s *Scraper) Scrape(filename string) (ScrapedData, error) {

	// Visit the list of characters page and gather URLs
	s.GetLocationURLs()
	s.GetCharacterURLs()

	// Start scraping each character
	s.ScrapeCharacters()
	s.ScrapeLocations()

	err := s.WriteDataToJSON(s.CharacterDataChannel, s.LocationDataChannel)
	if err != nil {
		return ScrapedData{}, err
	}

	return ScrapedData{
		CharacterURLs:      s.CharacterURLs,
		AmountOfCharacters: len(s.CharacterURLs),
	}, nil

}

// GetCharacterURLs gathers all unique character URLs from the list page
func (s *Scraper) GetCharacterURLs() {
	var scrapedUrls []string
	s.Logger.Infoln("Started character scraper...")

	printDebug("Getting character URLs...", s)

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	c.OnHTML("div#portal_frame a[title]", func(e *colly.HTMLElement) {
		characterURL := e.Request.AbsoluteURL(e.Attr("href"))
		if _, exists := s.URLSet[characterURL]; !exists {
			s.URLSet[characterURL] = struct{}{}
			printDebug("Found character URL: "+characterURL, s)
			scrapedUrls = append(scrapedUrls, characterURL)
			s.CharacterURLs = append(s.CharacterURLs, characterURL)
		}
	})

	c.Visit(fmt.Sprintf("%s/wiki/List_of_Characters", BaseFandomURL))

}

func (s *Scraper) GetLocationURLs() {
	locationUrls := []string{}
	s.Logger.Infoln("Started location scraper...")

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	// List of valid major locations to track
	validMajorLocations := []string{"Central Lands", "Northern Lands", "Southern Lands"}

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		// Get the major location name
		locationName := e.DOM.Find("span.mw-headline").Text()

		if slices.Contains(validMajorLocations, locationName) {
			fmt.Println(locationName)
			e.DOM.NextAllFiltered("ul").Find("li a").Each(func(i int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists && !(strings.Contains(href, ":") || strings.Contains(href, "cite")) {
					locationUrls = append(locationUrls, BaseFandomURL+href)
				}
			})
		}

	})

	// Start scraping
	c.Visit(fmt.Sprintf("%s/wiki/Locations", BaseFandomURL))

	locationUrls = removeDuplicates(locationUrls)
	fmt.Printf("Found %d locations\n", len(locationUrls))

	// Now locationUrls should have all the collected URLs
	s.LocationURLs = locationUrls
}

// ScrapeCharacters starts the scraping process for each character URL
func (s *Scraper) ScrapeCharacters() {
	wg := &sync.WaitGroup{}

	for _, url := range s.CharacterURLs {
		wg.Add(1)
		printDebug("Scraping character: "+url, s)
		go scrapeCharacter(url, wg, s.CharacterDataChannel)
	}

	go func() {
		wg.Wait()
		close(s.CharacterDataChannel)
	}()
}

func (s *Scraper) ScrapeLocations() {
	wg := &sync.WaitGroup{}

	for _, url := range s.LocationURLs {
		wg.Add(1)
		// printDebug("Scraping location: "+url, s)
		go scrapeLocation(url, wg, s.LocationDataChannel)
	}

	go func() {
		wg.Wait()
		close(s.LocationDataChannel)
	}()

}

func (s *Scraper) SetDebug(shouldDebug bool) {
	s.ShouldDebug = shouldDebug
}

func (s *Scraper) WriteDataToJSON(characterChannel chan *models.Character, locationChannel chan *models.Location) error {
	file, err := os.Create("characters.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	var characters []*models.Character
	var locations []string

	// Read from DataChannel until it's closed
	for data := range characterChannel {
		characters = append(characters, data)
	}

	for data := range locationChannel {
		locations = append(locations, data.URL)
	}

	// Encode all characters to JSON
	err = encoder.Encode(characters)
	if err != nil {
		return err
	}

	fmt.Println(locations)

	storage.CharactersData, err = data.LoadCharacters("characters.json")
	if err != nil {
		return err
	}

	return nil
}

// func newLocationURLs() *Locations {
// 	return &Locations{
// 		Central: []string{},
// 		Nothern: &NothernLocations{
// 			NothernPlateau:    []string{},
// 			ImperialTerritory: []string{},
// 			Ende:              []string{},
// 		},
// 		Southern: []string{},
// 	}
// }

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
