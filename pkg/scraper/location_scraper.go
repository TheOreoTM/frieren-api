package scraper

import (
	"fmt"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/theoreotm/frieren-api/models"
)

func scrapeLocation(url string, wg *sync.WaitGroup, channel chan *models.Character) {
	defer wg.Done()
	fmt.Println("Scraping location: " + url)

	channel <- models.NewCharacter(url)

}

func extractLocations(e *goquery.Selection) map[string]string {
	locations := make(map[string]string)

	for next := e.Parent().Next(); next.Length() > 0; next = next.Next() {
		if next.Is("h2") { // Stop if a new heading is encountered
			break
		}

		if next.Is("p") {
			locations["default"] = cleanText(next)
		}

		if !next.Is("ul") { // Stop if we encounter a figure element (aka: ability shown in a picture)
			continue
		}

		if next.Is("ul") {
			next.Contents().Each(func(i int, s *goquery.Selection) {
				location, url := getLocation(s)
				fmt.Printf("location: %q, url: %s\n", location, url)
				locations[location] = url
			})
		}

	}

	return locations
}

func getLocation(selection *goquery.Selection) (string, string) {
	location := ""
	url := ""

	selection.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("a") {
			href, exists := s.Attr("href")
			if exists {
				url = href
			}
		}

		location = s.Text()
	})

	return location, url
}
