package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

type Scraper struct {
	CharacterURLs []string
	URLSet        map[string]struct{}
	DataChannel   chan Character
	ShouldDebug   bool
}

// NewScraper initializes a new Scraper
func NewScraper() *Scraper {
	return &Scraper{
		CharacterURLs: []string{},
		URLSet:        make(map[string]struct{}),
		DataChannel:   make(chan Character),
		ShouldDebug:   false,
	}
}

// GetCharacterURLs gathers all unique character URLs from the list page
func (s *Scraper) GetCharacterURLs(wg *sync.WaitGroup, ctx *Ctx) {
	fmt.Println("Scraper started...")

	debug("Getting character URLs...", s)

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	c.OnHTML("div#portal_frame a[title]", func(e *colly.HTMLElement) {
		characterURL := e.Request.AbsoluteURL(e.Attr("href"))
		if _, exists := s.URLSet[characterURL]; !exists {
			s.URLSet[characterURL] = struct{}{}
			debug("Found character URL: "+characterURL, s)
			s.CharacterURLs = append(s.CharacterURLs, characterURL)
		}
	})

	c.Visit("https://frieren.fandom.com/wiki/List_of_Characters")
}

// ScrapeCharacters starts the scraping process for each character URL
func (s *Scraper) ScrapeCharacters(wg *sync.WaitGroup, ctx *Ctx) {
	for _, url := range s.CharacterURLs {
		wg.Add(1)
		debug("Scraping character: "+url, s)
		go scrapeCharacter(wg, ctx)
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
	var characters []Character

	// Read from DataChannel until it's closed
	for data := range s.DataChannel {
		characters = append(characters, data)
	}

	// Encode all characters to JSON
	err = encoder.Encode(characters)
	if err != nil {
		return err
	}

	return nil
}

// WriteDataToCSV writes the scraped data to a CSV file
func (s *Scraper) WriteDataToCSV(filename string) error {
	// Create or open the CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	err = writer.Write([]string{"URL", "Character", "Class", "Gender", "Rank", "Species"})
	if err != nil {
		return err
	}

	// Write data for each character
	for data := range s.DataChannel {
		row := []string{
			data.URL,
			data.Data["character"],
			data.Data["class"],
			data.Data["gender"],
			data.Data["rank"],
			data.Data["species"],
		}
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
