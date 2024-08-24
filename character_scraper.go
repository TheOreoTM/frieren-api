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
	URL       string            `json:"url"`
	Data      map[string]string `json:"data"`
	Abilities Abilities         `json:"abilities"`
}

func scrapeCharacter(url string, wg *sync.WaitGroup, channel chan Character) {
	defer wg.Done()

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	data := Character{URL: url, Data: make(map[string]string)}

	c.OnHTML(".mw-page-title-main", func(e *colly.HTMLElement) {
		data.Data["character"] = cleanText(e.DOM)
	})

	getCharInfo("species", data, c)
	getCharInfo("gender", data, c)
	getCharInfo("class", data, c)
	getCharInfo("rank", data, c)

	// Extract abilities and store them in the data struct
	c.OnHTML("h2 span#Abilities", func(e *colly.HTMLElement) {
		data.Abilities = extractAbilities(e.DOM)
	})

	c.Visit(data.URL)
	channel <- data
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

func extractAbilities(e *goquery.Selection) map[string]string {
	abilities := make(map[string]string)
	abilities["default"] = ""

	// Find the "Abilities" section
	// Traverse the siblings of the "Abilities" heading
	for next := e.Parent().Next(); next.Length() > 0; next = next.Next() {
		if next.Is("h2") { // Stop if a new heading is encountered
			break
		}

		if !next.Is("ul") { // Stop if we encounter a figure element (aka: ability shown in a picture)
			continue
		}

		flattenedList := flattenList(next)
		for _, ability := range flattenedList {
			sections := strings.Split(ability, ": ")
			name := strings.TrimSpace(strings.Join(sections[:len(sections)-1], ": "))
			description := strings.TrimSpace(strings.Replace(ability, name+":", "", 1))

			if name != "" && description != "" {
				abilities[name] = description
			}
		}

		// abiltyText := next.Text()
		// listOfAbilties := strings.Split(abiltyText, "\n")

		// for _, ability := range listOfAbilties {
		// 	fmt.Printf("Found Ability: %q\n", ability)
		// }

		// 	if next.Is("ul") {
		// 		processItemLists(next, abilities)
		// 	}

		// 	if next.Is("p") { // Capture text paragraphs within the section
		// 		if abilities["default"] == "" {
		// 			abilities["default"] = cleanText(next)
		// 		} else {
		// 			abilities["default"] += "\n" + cleanText(next)
		// 		}
		// 	}
	}

	return abilities
}

// func processItemLists(ul *goquery.Selection, abilities map[string]string) {
// 	ul.Children().Each(func(_ int, li *goquery.Selection) {
// 		if li.Is("li") {
// 			if li.Find("ul").Length() > 0 {
// 				// We do this because if the ability is a parent ability there is a nested <ul> element in the <li>
// 				// We need to extract the parent ability from the <li> and the nested <ul> from the <li>
// 				// Otherwise the entire <ul> is also extracted as an ability description
// 				processNestedAbilities(li, abilities)

// 				// Now process the nested <ul> abilities
// 				li.Find("ul li").Each(func(_ int, nestedLi *goquery.Selection) {
// 					// processNestedAbilities(nestedLi, abilities)
// 					nestedText := cleanText(nestedLi)
// 					nestedSections := strings.Split(nestedText, ": ")
// 					nestedName := strings.TrimSpace(strings.Join(nestedSections[:len(nestedSections)-1], ": "))
// 					nestedDescription := strings.TrimSpace(strings.Replace(nestedText, nestedName+":", "", 1))

// 					if nestedName != "" && nestedDescription != "" {
// 						abilities[nestedName] = nestedDescription
// 					}
// 				})

// 				// Now process the nested nested <ul> abilities, yeah this is a bit of a mess but ill fix it later
// 				li.Find("ul li ul li").Each(func(_ int, nestedNestedLi *goquery.Selection) {
// 					// processNestedAbilities(nestedNestedLi, abilities)

// 					nestedNestedText := cleanText(nestedNestedLi)
// 					nestedNestedSections := strings.Split(nestedNestedText, ": ")
// 					nestedNestedName := strings.TrimSpace(strings.Join(nestedNestedSections[:len(nestedNestedSections)-1], ": "))
// 					nestedNestedDescription := strings.TrimSpace(strings.Replace(nestedNestedText, nestedNestedName+":", "", 1))

// 					if nestedNestedName != "" && nestedNestedDescription != "" {
// 						abilities[nestedNestedName] = nestedNestedDescription
// 					}
// 				})
// 			} else {
// 				// Extract parent ability
// 				fullText := cleanText(li)
// 				sections := strings.Split(fullText, ": ")
// 				name := strings.TrimSpace(strings.Join(sections[:len(sections)-1], ": "))
// 				description := strings.TrimSpace(strings.Replace(fullText, name+":", "", 1))

// 				if name != "" && description != "" {
// 					abilities[name] = description
// 				}
// 			}
// 		}
// 	})
// }

// func processNestedAbilities(li *goquery.Selection, abilities map[string]string) {
// 	var parentFullText strings.Builder

// 	// Here we are adding text into the parentFullText builder until we reach the nested <ul> element
// 	// This is done by iterating over the <li> contents and appending text to the parentFullText builder
// 	li.Contents().Each(func(i int, s *goquery.Selection) {
// 		if s.Is("ul") {
// 			return
// 		}

// 		if nodeText := cleanText(s); nodeText != "" {
// 			parentFullText.WriteString(nodeText)
// 		}
// 	})

// 	parentSections := strings.Split(parentFullText.String(), ": ")
// 	parentName := strings.TrimSpace(strings.Join(parentSections[:len(parentSections)-1], ": "))
// 	parentDescription := strings.TrimSpace(strings.Replace(parentFullText.String(), parentName+":", "", 1))

// 	if parentName != "" && parentDescription != "" {
// 		abilities[parentName] = parentDescription
// 	}
// }

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
