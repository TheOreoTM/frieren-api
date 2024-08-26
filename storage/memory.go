package storage

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/theoreotm/frieren-api/models"
	"github.com/theoreotm/frieren-api/pkg/data"
)

// Initialize the data when the application starts
func init() {
	var err error
	characters, err = data.LoadCharacters("characters.json")
	if err != nil {
		logrus.Panicf("Failed to load character data: %v", err)
	}
}

type MemoryStorage struct {
}

// Global variable to hold character data
var characters *models.Characters

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) GetCharacters() models.Characters {
	return *characters
}

func (s *MemoryStorage) GetCharacter(name string) (models.Character, error) {
	for _, character := range characters.Characters {
		if strings.EqualFold(character.Data["character"], name) {
			return character, nil
		}
	}

	return models.Character{}, errors.New("character not found")
}
