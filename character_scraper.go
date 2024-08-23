package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Abilities map[string]string

type Character struct {
	URL       string
	Data      map[string]string
	Abilities Abilities
}

func scrapeCharacter(url string, wg *sync.WaitGroup, ch chan<- Character) {
	defer wg.Done()

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	data := Character{URL: url, Data: make(map[string]string)}
	data.URL = url

	c.OnHTML(".mw-page-title-main", func(e *colly.HTMLElement) {
		data.Data["character"] = cleanText(e.DOM)
	})

	getCharInfo("species", data, c)
	getCharInfo("gender", data, c)
	getCharInfo("class", data, c)
	getCharInfo("rank", data, c)

	// Extract abilities and store them in the data struct
	c.OnHTML("div#content", func(e *colly.HTMLElement) {
		data.Abilities = extractAbilities(e)
	})

	c.Visit(url)
	ch <- data
}

func getCharInfo(info string, character Character, c *colly.Collector) {
	c.OnHTML(fmt.Sprintf("div[data-source='%s'] .pi-data-value", info), func(e *colly.HTMLElement) {
		text := cleanText(e.DOM)
		if text == "" {
			return
		}
		character.Data[info] = text
	})
}

func extractAbilities(e *colly.HTMLElement) map[string]string {
	abilities := make(map[string]string)
	abilities["default"] = ""

	// Find the "Abilities" section
	e.DOM.Find("h2 span#Abilities").Each(func(_ int, s *goquery.Selection) {
		// Traverse the siblings of the "Abilities" heading
		for next := s.Parent().Next(); next.Length() > 0; next = next.Next() {
			if next.Is("h2") { // Stop if a new heading is encountered
				break
			}

			if next.Is("ul") {
				processItemLists(next, abilities)
			}

			if next.Is("p") { // Capture text paragraphs within the section
				if abilities["default"] == "" {
					abilities["default"] = cleanText(next)
				} else {
					abilities["default"] += "\n" + cleanText(next)
				}
			}
		}
	})

	return abilities
}

func processItemLists(ul *goquery.Selection, abilities map[string]string) {
	ul.Children().Each(func(_ int, li *goquery.Selection) {
		if li.Is("li") {
			if li.Find("ul").Length() > 0 {
				// We do this because if the ability is a parent ability there is a nested <ul> element in the <li>
				// We need to extract the parent ability from the <li> and the nested <ul> from the <li>
				// Otherwise the entire <ul> is also extracted as an ability description
				processNestedAbilities(li, abilities)

				// Now process the nested <ul> abilities
				li.Find("ul li").Each(func(_ int, nestedLi *goquery.Selection) {
					processNestedAbilities(nestedLi, abilities)
					// nestedText := cleanText(nestedLi)
					// nestedSections := strings.Split(nestedText, ": ")
					// nestedName := strings.TrimSpace(strings.Join(nestedSections[:len(nestedSections)-1], ": "))
					// nestedDescription := strings.TrimSpace(strings.Replace(nestedText, nestedName+":", "", 1))

					// if nestedName != "" && nestedDescription != "" {
					// 	abilities[nestedName] = nestedDescription
					// }
				})

				// Now process the nested nested <ul> abilities, yeah this is a bit of a mess but ill fix it later
				li.Find("ul li ul li").Each(func(_ int, nestedNestedLi *goquery.Selection) {
					processNestedAbilities(nestedNestedLi, abilities)

					// nestedNestedText := cleanText(nestedNestedLi)
					// nestedNestedSections := strings.Split(nestedNestedText, ": ")
					// nestedNestedName := strings.TrimSpace(strings.Join(nestedNestedSections[:len(nestedNestedSections)-1], ": "))
					// nestedNestedDescription := strings.TrimSpace(strings.Replace(nestedNestedText, nestedNestedName+":", "", 1))

					// if nestedNestedName != "" && nestedNestedDescription != "" {
					// 	abilities[nestedNestedName] = nestedNestedDescription
					// }
				})
			} else {
				// Extract parent ability
				fullText := cleanText(li)
				sections := strings.Split(fullText, ": ")
				name := strings.TrimSpace(strings.Join(sections[:len(sections)-1], ": "))
				description := strings.TrimSpace(strings.Replace(fullText, name+":", "", 1))

				if name != "" && description != "" {
					abilities[name] = description
				}
			}
		}
	})
}

func processNestedAbilities(li *goquery.Selection, abilities map[string]string) {
	var parentFullText strings.Builder

	// Here we are adding text into the parentFullText builder until we reach the nested <ul> element
	// This is done by iterating over the <li> contents and appending text to the parentFullText builder
	li.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("ul") {
			return
		}

		if nodeText := cleanText(s); nodeText != "" {
			parentFullText.WriteString(nodeText)
		}
	})

	parentSections := strings.Split(parentFullText.String(), ": ")
	parentName := strings.TrimSpace(strings.Join(parentSections[:len(parentSections)-1], ": "))
	parentDescription := strings.TrimSpace(strings.Replace(parentFullText.String(), parentName+":", "", 1))

	if parentName != "" && parentDescription != "" {
		abilities[parentName] = parentDescription
	}

}
