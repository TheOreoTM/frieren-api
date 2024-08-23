package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	var wg sync.WaitGroup
	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))
	c.Limit(&colly.LimitRule{Parallelism: 10, Delay: 500 * time.Millisecond})

	// Initialize the scraper and URL list
	scraper := NewScraper()

	// Visit the list of characters page and gather URLs
	scraper.GetCharacterURLs(&wg)

	// Start scraping each character
	scraper.ScrapeCharacters(&wg, c)

	// // Scrape one character
	// scrapeCharacter("https://frieren.fandom.com/wiki/Sein", &wg, scraper.DataChannel)

	// Wait for all scraping goroutines to finish
	go func() {
		wg.Wait()
		close(scraper.DataChannel)
	}()

	// err := scraper.WriteDataToCSV("characters.csv")
	// if err != nil {
	// 	fmt.Println("Error writing data to CSV:", err)
	// }
	err := scraper.WriteDataToJSON("characters.json")
	if err != nil {
		fmt.Println("Error writing data to JSON:", err)
	}

	fmt.Printf("Scraped %d characters\n", len(scraper.CharacterURLs))
}
