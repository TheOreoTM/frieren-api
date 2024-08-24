package main

type Context struct {
	CurrentURL string
	Scraper    *Scraper
}

func NewContext() *Context {
	return &Context{
		Scraper: NewScraper(),
	}
}

func (c *Context) SetCurrentURL(url string) {
	c.CurrentURL = url
}
