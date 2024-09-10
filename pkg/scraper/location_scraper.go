package scraper

import (
	"sync"

	"github.com/theoreotm/frieren-api/models"
)

func scrapeLocation(url string, wg *sync.WaitGroup, channel chan *models.Location) {
	defer wg.Done()
	location := models.NewLocation(url)

	

	channel <- location
}
