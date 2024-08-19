package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Initialize the scraper and URL list
	scraper := NewScraper()

	// Visit the list of characters page and gather URLs
	scraper.GetCharacterURLs(&wg)

	// Start scraping each character
	scraper.ScrapeCharacters(&wg)

	// Wait for all scraping goroutines to finish
	go func() {
		wg.Wait()
		close(scraper.DataChannel)
	}()

	// Print the collected data as JSON
	for data := range scraper.DataChannel {
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling data:", err)
			continue
		}

		fmt.Printf("Character: %s\n", data["character"])
		fmt.Println(string(jsonData) + "\n")
	}
}
