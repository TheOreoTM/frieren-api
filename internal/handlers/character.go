package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/theoreotm/frieren-api/models"
	"github.com/theoreotm/frieren-api/pkg/data"
)

// Global variable to hold character data
var characters *models.Characters

// Initialize the data when the application starts
func init() {
	var err error
	characters, err = data.LoadCharacters("characters.json")
	if err != nil {
		logrus.Panicf("Failed to load character data: %v", err)
	}
}

// GetCharacter handles the GET /character/{name} endpoint.
func GetCharacter(w http.ResponseWriter, r *http.Request, logger *logrus.Logger) {
	vars := mux.Vars(r)
	name := vars["name"]

	for _, character := range characters.Characters {
		if strings.EqualFold(character.Data["character"], name) {
			logger.Infof("Character found: %s", character.Data["character"])
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(character)
			return
		}
	}

	logger.Warnf("Character not found: %s", name)
	http.Error(w, "Character not found", http.StatusNotFound)
}

// GetCharacters handles the GET /characters endpoint.
func GetCharacters(w http.ResponseWriter, r *http.Request, logger *logrus.Logger) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
