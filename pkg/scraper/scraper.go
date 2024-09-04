package scraper

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/theoreotm/frieren-api/models"
	"github.com/theoreotm/frieren-api/pkg/data"
	"github.com/theoreotm/frieren-api/storage"
)

type Scraper struct {
	CharacterURLs []string
	URLSet        map[string]struct{}
	DataChannel   chan *models.Character
	ShouldDebug   bool
	Logger        *logrus.Logger
}

type ScrapedData struct {
	CharacterURLs      []string
	AmountOfCharacters int
}

// NewScraper initializes a new Scraper
func NewScraper(shouldDebug bool, logger *logrus.Logger) *Scraper {
	return &Scraper{
		CharacterURLs: []string{},
		URLSet:        make(map[string]struct{}),
		DataChannel:   make(chan *models.Character),
		ShouldDebug:   shouldDebug,
		Logger:        logger,
	}
}

func (s *Scraper) Scrape(filename string) (ScrapedData, error) {
	var wg sync.WaitGroup

	// Visit the list of characters page and gather URLs
	scrapedUrls := s.GetCharacterURLs(&wg)

	// Start scraping each character
	s.ScrapeCharacters(&wg)

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
func (s *Scraper) GetCharacterURLs(wg *sync.WaitGroup) []string {
	var scrapedUrls []string
	s.Logger.Infoln("Started scraper...")

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

// ScrapeCharacters starts the scraping process for each character URL
func (s *Scraper) ScrapeCharacters(wg *sync.WaitGroup) {
	for _, url := range s.CharacterURLs {
		wg.Add(1)
		printDebug("Scraping character: "+url, s)
		go scrapeCharacter(url, wg, s.DataChannel)
	}
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
