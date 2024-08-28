package storage

import (
	"errors"
	"strings"

	"github.com/theoreotm/frieren-api/models"
)

type MemoryStorage struct {
}

// Global variable to hold character data
var CharactersData *models.Characters

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) GetCharacters() models.Characters {
	return *CharactersData
}

func (s *MemoryStorage) GetCharacter(name string) (models.Character, error) {
	for _, character := range CharactersData.Characters {
		if strings.EqualFold(character.Data["character"], name) {
			return character, nil
		}
	}

	return models.Character{}, errors.New("character not found")
}
