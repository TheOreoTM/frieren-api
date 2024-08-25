package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	// Import your scraper package if needed
)

// GetCharacters handles the GET /characters endpoint.
func GetCharacters(w http.ResponseWriter, r *http.Request) {
	// Example response
	response := map[string]string{"message": "List of characters"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCharacter handles the GET /character/{id} endpoint.
func GetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	
	// Example response
	response := map[string]string{"id": name, "name": "Character Name"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
