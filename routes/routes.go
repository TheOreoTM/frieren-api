package routes

import (
	"github.com/gorilla/mux"
	"github.com/theoreotm/frieren-api/internal/handlers"
)

// SetupRoutes defines the API routes.
func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/characters", handlers.GetCharacters).Methods("GET")
	r.HandleFunc("/characters/{name}", handlers.GetCharacter).Methods("GET")
	// Add more routes as needed
}
