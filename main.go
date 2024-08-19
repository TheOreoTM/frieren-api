package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func scrape(url string, wg *sync.WaitGroup, ch chan<- map[string]string) {
	defer wg.Done()

	c := colly.NewCollector()

	data := make(map[string]string)

	c.OnHTML(".mw-page-title-main", func(e *colly.HTMLElement) {
		data["character"] = e.Text
	})

	c.OnHTML("div[data-source='species'] .pi-data-value", func(e *colly.HTMLElement) {
		species := e.ChildText("a")
		data["species"] = species
	})

	c.OnHTML("div[data-source='gender'] .pi-data-value", func(e *colly.HTMLElement) {
		gender := e.ChildText("a")
		data["gender"] = gender
	})

	c.OnHTML("div[data-source='class'] .pi-data-value", func(e *colly.HTMLElement) {
		class := e.ChildText("a")
		data["class"] = class
	})

	c.OnHTML("div[data-source='rank'] .pi-data-value", func(e *colly.HTMLElement) {
		rank := cleanText(e.DOM)
		data["rank"] = rank
	})

	// Visit the page and once done, send the data through the channel
	c.Visit(url)
	ch <- data
}

// Helper function to clean up text, removing unwanted tags like <sup> or <br>
func cleanText(selection *goquery.Selection) string {
	var rankText strings.Builder

	// Loop over each child element and append only text nodes, ignoring <sup>, <br>, etc.
	selection.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("sup") {
			// Ignore <sup> elements
			return
		}
		if s.Is("br") {
			// Replace <br> elements with a newline character
			rankText.WriteString("\n")
			return
		}
		if nodeText := s.Text(); nodeText != "" {
			rankText.WriteString(strings.TrimSpace(nodeText) + " ")
		}
	})

	return strings.TrimSpace(rankText.String())
}

func main() {
	var wg sync.WaitGroup
	urls := []string{
		"https://frieren.fandom.com/wiki/Frieren",
		"https://frieren.fandom.com/wiki/Himmel",
		"https://frieren.fandom.com/wiki/Heiter",
	}

	// Create a channel to collect the scraped data
	dataChannel := make(chan map[string]string, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go scrape(url, &wg, dataChannel)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(dataChannel)
	}()

	for data := range dataChannel {
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling data:", err)
			continue
		}

		// Print the collected data as JSON
		fmt.Println(string(jsonData) + "\n")
	}
}
