package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/theoreotm/frieren-api/models"
	"github.com/theoreotm/frieren-api/pkg/data"
	"github.com/theoreotm/frieren-api/storage"
)

type Scraper struct {
	CharacterURLs []string
	LocationURLs  *LocationURLs
	URLSet        map[string]struct{}
	DataChannel   chan *models.Character
	ShouldDebug   bool
	Logger        *logrus.Logger
}

type ScrapedData struct {
	CharacterURLs      []string
	AmountOfCharacters int
}

type LocationURLs struct {
	Central  []string
	Nothern  *NothernLocationURLs
	Southern []string
}

type NothernLocationURLs struct {
	NothernPlateau    []string
	ImperialTerritory []string
	Ende              []string
}

const (
	BaseFandomURL = "https://frieren.fandom.com"
)

// NewScraper initializes a new Scraper
func NewScraper(shouldDebug bool, logger *logrus.Logger) *Scraper {
	return &Scraper{
		CharacterURLs: []string{},
		LocationURLs:  newLocationURLs(),
		URLSet:        make(map[string]struct{}),
		DataChannel:   make(chan *models.Character),
		ShouldDebug:   shouldDebug,
		Logger:        logger,
	}
}

func (s *Scraper) Scrape(filename string) (ScrapedData, error) {
	var wg sync.WaitGroup

	// Visit the list of characters page and gather URLs
	s.GetLocationURLs()
	scrapedUrls := s.GetCharacterURLs()

	// Start scraping each character

	s.ScrapeCharacters(&wg)
	s.ScrapeLocations(&wg)

	// Wait for all scraping goroutines to finish
	go func() {
		wg.Wait()
		close(s.DataChannel)
	}()

	err := s.WriteDataToJSON(filename)
	if err != nil {
		return ScrapedData{}, err
	}

	return ScrapedData{
		CharacterURLs:      scrapedUrls,
		AmountOfCharacters: len(scrapedUrls),
	}, nil

}

// GetCharacterURLs gathers all unique character URLs from the list page
func (s *Scraper) GetCharacterURLs() []string {
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

	c.Visit("https://frieren.fandom.com/wiki/List_of_Characters")

	return scrapedUrls
}

func (s *Scraper) GetLocationURLs() {
	locationUrls := &LocationURLs{}
	s.Logger.Infoln("Started location scraper...")

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	validMajorLocations := []string{"Central Lands", "Northern Lands", "Southern Lands"}

	c.OnHTML("h2 span#Central_Lands", func(e *colly.HTMLElement) {
		section := cleanText(e.DOM)
		if !slices.Contains(validMajorLocations, section) {
			return
		}

		locations := extractLocations(e.DOM)
		fmt.Printf("Found %d locations in %s\n", len(locations), section)

	})

	c.Visit("https://frieren.fandom.com/wiki/Locations")

	fmt.Println(locationUrls)

}

// ScrapeCharacters starts the scraping process for each character URL
func (s *Scraper) ScrapeCharacters(wg *sync.WaitGroup) {
	for _, url := range s.CharacterURLs {
		wg.Add(1)
		printDebug("Scraping character: "+url, s)
		go scrapeCharacter(url, wg, s.DataChannel)
	}
}

func (s *Scraper) ScrapeLocations(wg *sync.WaitGroup) {
	for _, url := range s.LocationURLs.Central {
		wg.Add(1)
		printDebug("Scraping location: "+url, s)
		go scrapeLocation(url, wg, s.DataChannel)
	}

}

func (s *Scraper) SetDebug(shouldDebug bool) {
	s.ShouldDebug = shouldDebug
}

func (s *Scraper) WriteDataToJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// Create a slice to hold all the characters
	var characters []*models.Character

	// Read from DataChannel until it's closed
	for data := range s.DataChannel {
		characters = append(characters, data)
	}

	// Encode all characters to JSON
	err = encoder.Encode(characters)
	if err != nil {
		return err
	}

	storage.CharactersData, err = data.LoadCharacters(filename)
	if err != nil {
		return err
	}

	return nil
}

func newLocationURLs() *LocationURLs {
	return &LocationURLs{
		Central: []string{},
		Nothern: &NothernLocationURLs{
			NothernPlateau:    []string{},
			ImperialTerritory: []string{},
			Ende:              []string{},
		},
		Southern: []string{},
	}
}
