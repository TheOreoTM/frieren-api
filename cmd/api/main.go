package main

import (
	"net/http"

	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	"github.com/theoreotm/frieren-api/api"
	"github.com/theoreotm/frieren-api/config"
	"github.com/theoreotm/frieren-api/pkg/scraper"
	"github.com/theoreotm/frieren-api/storage"
)

func main() {
	cfg := config.LoadConfig()
	listenAddr := ":" + cfg.Port
	logger := config.SetupLogger(cfg.LogLevel)
	scraper := scraper.NewScraper(false, logger)

	data, err := scraper.Scrape("characters.json")

	r := mux.NewRouter()
	r.Use(muxlogrus.NewLogger().Middleware)

	store := storage.NewMemoryStorage()

	server := api.NewServer(listenAddr, store, logger)
	server.Start(r, logger)

	if err != nil {
		logger.Errorf("Error scraping: %v", err)
	}

	logger.Infof("Scraped %d characters", data.AmountOfCharacters)

	logger.Infof("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}

}
