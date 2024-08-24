package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Helper function to clean up text, removing unwanted tags like <br> or <sup>
func cleanText(selection *goquery.Selection, ctx *Ctx) string {
	var cleanedText strings.Builder

	selection.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is(".reference") {
			fmt.Printf("Found sup: %s\n\n\n\n", s.Text())
			href, exists := s.Find("a").Attr("href")
			if exists {
				cleanedText.WriteString(fmt.Sprintf("[%s](%s)", strings.TrimSpace(s.Text()), fmt.Sprintf("%s%s", ctx.CurrentURL, href)))
			}
		}

		if s.Is("br") {
			cleanedText.WriteString("\n") // Handle <br> as newlines
			return
		}
		if s.Is("a") {
			href, exists := s.Attr("href")
			if exists {
				// Format hyperlink as Markdown: [text](link)
				cleanedText.WriteString(fmt.Sprintf("[%s](%s)", strings.TrimSpace(s.Text()), fmt.Sprintf("https://frieren.fandom.com%s", href)))
			}
			return
		}
		if nodeText := s.Text(); nodeText != "" {
			cleanedText.WriteString(nodeText)
		}
	})

	return strings.TrimSpace(cleanedText.String())
}

func debug(msg string, s *Scraper) {
	if s.ShouldDebug {
		fmt.Println(msg)
	}
}
