package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

type Scraper struct {
	CharacterURLs []string
	URLSet        map[string]struct{}
	DataChannel   chan map[string]string
	ShouldDebug   bool
}

// NewScraper initializes a new Scraper
func NewScraper() *Scraper {
	return &Scraper{
		CharacterURLs: []string{},
		URLSet:        make(map[string]struct{}),
		DataChannel:   make(chan map[string]string),
		ShouldDebug:   false,
	}
}

// GetCharacterURLs gathers all unique character URLs from the list page
func (s *Scraper) GetCharacterURLs(wg *sync.WaitGroup) {
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
func (s *Scraper) ScrapeCharacters(wg *sync.WaitGroup) {
	for _, url := range s.CharacterURLs {
		wg.Add(1)
		debug("Scraping character: "+url, s)
		go scrapeCharacter(url, wg, s.DataChannel)
	}
}
