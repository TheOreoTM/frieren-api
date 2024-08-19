package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

func scrapeCharacter(url string, wg *sync.WaitGroup, ch chan<- map[string]string) {
	defer wg.Done()

	c := colly.NewCollector(colly.AllowedDomains("frieren.fandom.com"))

	data := make(map[string]string)
	data["url"] = url

	c.OnHTML(".mw-page-title-main", func(e *colly.HTMLElement) {
		data["character"] = cleanText(e.DOM)
	})

	getCharInfo("species", data, c)
	getCharInfo("gender", data, c)
	getCharInfo("class", data, c)
	getCharInfo("rank", data, c)

	c.Visit(url)
	ch <- data
}

func getCharInfo(info string, data map[string]string, c *colly.Collector) {
	c.OnHTML(fmt.Sprintf("div[data-source='%s'] .pi-data-value", info), func(e *colly.HTMLElement) {
		data[info] = cleanText(e.DOM)
	})
}
