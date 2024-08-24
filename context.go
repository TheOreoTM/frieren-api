package main

type Ctx struct {
	CurrentURL string
	Scraper    *Scraper
}

func NewContext() *Ctx {
	return &Ctx{
		Scraper: NewScraper(),
	}
}

func (c *Ctx) SetCurrentURL(url string) {
	c.CurrentURL = url
}
