package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theoreotm/frieren-api/config"
	"github.com/theoreotm/frieren-api/pkg/scraper"
	"github.com/theoreotm/frieren-api/routes"
)

func main() {
	cfg := config.LoadConfig()
	logger := config.SetupLogger(cfg.LogLevel)
	r := mux.NewRouter()
	routes.SetupRoutes(r, logger)
	scraper := scraper.NewScraper(false, logger)

	data, err := scraper.Scrape("characters.json")

	if err != nil {
		logger.Errorf("Error scraping: %v", err)
	}

	logger.Infof("Scraped %d characters", data.AmountOfCharacters)

	logger.Infof("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}

}
