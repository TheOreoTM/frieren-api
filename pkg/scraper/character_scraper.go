package scraper

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/theoreotm/frieren-api/models"
)

func scrapeCharacter(url string, wg *sync.WaitGroup, channel chan *models.Character) {
	defer wg.Done()

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	character := models.NewCharacter(url)

	// Scrape general character information
	c.OnHTML(".mw-page-title-main", func(e *colly.HTMLElement) {
		character.AddGeneralData("name", cleanText(e.DOM))
	})

	generalFields := []string{"species", "gender", "class", "rank", "age", "status", "affiliation", "relatives"}
	for _, field := range generalFields {
		// Capture the field within the loop
		field := field
		c.OnHTML(fmt.Sprintf("div[data-source='%s'] .pi-data-value", field), func(e *colly.HTMLElement) {
			data := cleanText(e.DOM)
			if data == "" {
				data = "Unknown"
			}
			character.AddGeneralData(field, data)
		})
	}

	seriesInformation := []string{"manga", "anime", "jpva", "enva"}
	for _, field := range seriesInformation {
		// Capture the field within the loop
		field := field
		c.OnHTML(fmt.Sprintf("div[data-source='%s'] .pi-data-value", field), func(e *colly.HTMLElement) {
			data := cleanText(e.DOM)
			if data == "" {
				data = "Unknown"
			}
			character.AddSeriesData(field, data)
		})
	}

	// Extract abilities and store them in the data struct
	c.OnHTML("h2 span#Abilities", func(e *colly.HTMLElement) {
		character.AddAbilities(extractAbilities(e.DOM))
	})

	c.Visit(character.URL)
	channel <- character
}

func getCharInfo(info string, defaultValue string, c *colly.Collector) string {
	var data string

	c.OnHTML(fmt.Sprintf("div[data-source='%s'] .pi-data-value", info), func(e *colly.HTMLElement) {
		text := cleanText(e.DOM)
		if text == "" {
			data = defaultValue
		}

		data = text
	})

	return data
}

func extractAbilities(e *goquery.Selection) map[string]string {
	abilities := make(map[string]string)

	// Find the "Abilities" section
	// Traverse the siblings of the "Abilities" heading
	for next := e.Parent().Next(); next.Length() > 0; next = next.Next() {
		if next.Is("h2") { // Stop if a new heading is encountered
			break
		}

		if next.Is("p") {
			if abilities["default"] == "" {
				abilities["default"] = cleanText(next)
			} else {
				abilities["default"] += "\n" + cleanText(next)
			}
		}

		if !next.Is("ul") { // Stop if we encounter a figure element (aka: ability shown in a picture)
			continue
		}

		flattenedList := flattenList(next)
		for _, ability := range flattenedList {

			if !strings.Contains(ability, ":") {
				abilities[ability] = "" // Stark edgecase: ability name has no description
				continue
			}

			sections := strings.Split(ability, ": ")
			name := strings.TrimSpace(strings.Join(sections[:len(sections)-1], ": "))
			description := strings.TrimSpace(strings.Replace(ability, name+":", "", 1))

			if name != "" || description != "" { // since starks ability list is just a list of ability names
				abilities[name] = description
			}
		}
	}

	return abilities
}

func flattenList(ul *goquery.Selection) []string {
	var result []string

	ul.Children().Each(func(_ int, li *goquery.Selection) {
		// Only process <li> elements that contain a <b> tag
		if li.Is("li") && li.Find("b").Length() > 0 {
			// Extract the parent <li> text excluding nested <ul> content
			nestedUl := li.Find("ul").First()
			var parentText string

			if nestedUl.Length() > 0 {
				// Extract text from the <li> without nested <ul>
				parentText = cleanText(li.Clone().ChildrenFiltered("ul").Remove().End())
			} else {
				parentText = cleanText(li)
			}

			if parentText != "" {
				result = append(result, parentText)
			}

			// Recursively process the nested <ul> if present
			if nestedUl.Length() > 0 {
				result = append(result, flattenList(nestedUl)...)
			}
		}
	})

	return result
}
