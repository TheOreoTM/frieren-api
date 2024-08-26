package storage

import "github.com/theoreotm/frieren-api/models"

type Storage interface {
	GetCharacters() models.Characters
	GetCharacter(name string) (models.Character, error)
}
