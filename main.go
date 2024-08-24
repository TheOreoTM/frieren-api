package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ctx := NewContext()
	ctx.SetCurrentURL("https://frieren.fandom.com/wiki/Frieren")

	// Initialize the scraper and URL list

	// Visit the list of characters page and gather URLs
	ctx.Scraper.GetCharacterURLs(&wg, ctx)

	// Start scraping each character
	// scraper.ScrapeCharacters(&wg)

	// Scrape one character
	scrapeCharacter(&wg, ctx)

	// Wait for all scraping goroutines to finish
	go func() {
		wg.Wait()
		close(ctx.Scraper.DataChannel)
	}()

	// err := scraper.WriteDataToCSV("characters.csv")
	// if err != nil {
	// 	fmt.Println("Error writing data to CSV:", err)
	// }
	err := ctx.Scraper.WriteDataToJSON("characters.json")
	if err != nil {
		fmt.Println("Error writing data to JSON:", err)
	}

	fmt.Printf("Scraped %d characters\n", len(ctx.Scraper.CharacterURLs))
}
