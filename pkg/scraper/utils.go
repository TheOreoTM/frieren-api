package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Helper function to clean up text, removing unwanted tags like <br> or <sup>
func cleanText(selection *goquery.Selection) string {
	var cleanedText strings.Builder

	selection.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is(".reference") {
			return
		}

		if s.Is("br") {
			cleanedText.WriteString("\n") // Handle <br> as newlines
			return
		}
		if s.Is("a") {
			href, exists := s.Attr("href")
			if exists {
				// Format hyperlink as Markdown: [text](link)
				hyperlink := ""

				if strings.HasPrefix(strings.TrimSpace(s.Text()), "/") {
					hyperlink = fmt.Sprintf("[%s](https://frieren.fandom.com%s)", strings.TrimSpace(s.Text()), href)
				} else {
					hyperlink = fmt.Sprintf("[%s](%s)", strings.TrimSpace(s.Text()), href)
				}
				cleanedText.WriteString(hyperlink)
			}
			return
		}

		if nodeText := s.Text(); nodeText != "" {
			cleanedText.WriteString(nodeText)
		}
	})

	return strings.TrimSpace(cleanedText.String())
}

func printDebug(msg string, s *Scraper) {
	if s.ShouldDebug {
		s.Logger.Debugln(msg)
	}
}
