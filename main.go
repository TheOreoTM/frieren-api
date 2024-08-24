package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	context := NewContext()
	context.SetCurrentURL("https://frieren.fandom.com/wiki/Frieren")

	// Initialize the scraper and URL list

	// Visit the list of characters page and gather URLs
	context.Scraper.GetCharacterURLs(&wg, context)

	// Start scraping each character
	// scraper.ScrapeCharacters(&wg)

	// Scrape one character
	scrapeCharacter(&wg, context)

	// Wait for all scraping goroutines to finish
	go func() {
		wg.Wait()
		close(context.Scraper.DataChannel)
	}()

	// err := scraper.WriteDataToCSV("characters.csv")
	// if err != nil {
	// 	fmt.Println("Error writing data to CSV:", err)
	// }
	err := context.Scraper.WriteDataToJSON("characters.json")
	if err != nil {
		fmt.Println("Error writing data to JSON:", err)
	}

	fmt.Printf("Scraped %d characters\n", len(context.Scraper.CharacterURLs))
}
