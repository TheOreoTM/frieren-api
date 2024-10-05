package api

import (
	"net/http"
	"strings"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/theoreotm/frieren-api/storage"
)

type Server struct {
	listenAddr string
	store      storage.Storage
	logger     *logrus.Logger
	*http.Server
}

func NewServer(listenAddr string, store storage.Storage, logger *logrus.Logger) *Server {
	httpServer := &http.Server{
		Addr: listenAddr,
	}

	return &Server{
		listenAddr: listenAddr,
		store:      store,
		logger:     logger,
		Server:     httpServer,
	}
}

func (s *Server) Start(r *mux.Router, logger *logrus.Logger) error {
	r.HandleFunc("/characters", makeHTTPHandleFunc(s.handleGetCharacters)).Methods("GET")
	r.HandleFunc("/characters/{name}", makeHTTPHandleFunc(s.handleGetCharacter)).Methods("GET")
	r.HandleFunc("/names", makeHTTPHandleFunc(s.handleGetNames)).Methods("GET")

	// Assign the router to the server's handler
	s.Handler = r

	// Start the server
	logger.Infof("Server starting on %s", s.listenAddr)
	return s.ListenAndServe()
}

// GetCharacter handles the GET /character/{name} endpoint.
func (s *Server) handleGetCharacter(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	name := vars["name"]

	character, err := s.store.GetCharacter(name)
	if err != nil {
		s.logger.Warnf("Character not found: %s", name)
		return WriteJSON(w, http.StatusNotFound, ApiError{Error: "character not found"})
	}

	return WriteJSON(w, http.StatusOK, character)
}

// GetCharacters handles the GET /characters endpoint.
func (s *Server) handleGetCharacters(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, s.store.GetCharacters())
}

func (s *Server) handleGetNames(w http.ResponseWriter, r *http.Request) error {
	// Get the names of all characters
	characterData := s.store.GetCharacters()
	names := make([]string, len(characterData.Characters))

	for i, character := range characterData.Characters {
		var charNames []string
		charNames = append(charNames, character.Data.Names.English)
		charNames = append(charNames, character.Data.Names.Japanese)
		charNames = append(charNames, character.Data.Names.Romaji)
		names[i] = strings.Join(charNames, ", ")
	}

	return WriteJSON(w, http.StatusOK, names)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
