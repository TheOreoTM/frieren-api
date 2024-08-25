package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/theoreotm/frieren-api/internal/handlers"
)

// SetupRoutes defines the API routes and attaches logging.
func SetupRoutes(r *mux.Router, logger *logrus.Logger) {
	r.HandleFunc("/characters", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCharacters(w, r, logger)
	}).Methods("GET")
	r.HandleFunc("/characters/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCharacter(w, r, logger)
	}).Methods("GET")
}
