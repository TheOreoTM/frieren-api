package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

type Character struct {
	URL  string
	Data map[string]string
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
