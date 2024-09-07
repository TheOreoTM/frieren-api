package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	scraper.SetDebug(true)

	// Scrape data
	data, err := scraper.Scrape("characters.json")
	if err != nil {
		logger.Errorf("Error scraping: %v", err)
		return
	}
	logger.Infof("Scraped %d characters", data.AmountOfCharacters)

	r := mux.NewRouter()
	r.Use(muxlogrus.NewLogger().Middleware)

	store := storage.NewMemoryStorage()

	// Create the server
	server := api.NewServer(listenAddr, store, logger)

	// Start the server in a goroutine
	go func() {
		if err := server.Start(r, logger); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Listen for OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	logger.Infof("Received signal %v, shutting down server", sig)

	// Graceful shutdown with a timeout context
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("HTTP shutdown error: %v", err)
	} else {
		logger.Info("Server gracefully stopped")
	}
}
