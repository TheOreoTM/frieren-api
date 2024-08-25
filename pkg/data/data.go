package data

import (
	"encoding/json"
	"os"

	"github.com/theoreotm/frieren-api/models"
)

// LoadCharacters loads character data from a JSON file.
func LoadCharacters(filename string) (*models.Characters, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var characters []models.Character
	if err := json.Unmarshal(file, &characters); err != nil {
		return nil, err
	}

	return &models.Characters{
		Characters: characters,
	}, nil
}
