package storage

import "github.com/theoreotm/frieren-api/models"

// Global variable to hold character data
var CharactersData *models.Characters = &models.Characters{}

type Storage interface {
	GetCharacters() models.Characters
	GetCharacter(name string) (models.Character, error)
}
