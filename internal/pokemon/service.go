package pokemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pokedex-clone/internal/storage"

	"github.com/gorilla/mux"
)

type Service struct {
	StorageAPI *storage.Store
}

func NewService(storage *storage.Store) *Service {
	return &Service{
		StorageAPI: storage,
	}
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	// var getPokemonRequest GetPokemonRequest
	// if err := json.NewDecoder(r.Body).Decode(&getPokemonRequest); err != nil {
	// 	s.renderError(w, "failed to decode request body", http.StatusBadRequest)
	// 	return
	// }

	if !s.StorageAPI.Exist(name) {
		s.renderError(w, fmt.Sprintf("missing pokemon with name %v", name), http.StatusNotFound)
		return
	}

	pokemon, ok := s.StorageAPI.Load(name)
	if !ok {
		s.renderError(w, fmt.Sprintf("failed to get pokemon %v", name), http.StatusBadRequest)
		return
	}

	p, ok := pokemon.(Pokemon)
	if !ok {
		s.renderError(w, fmt.Sprintf("unexpected object returned for %v, %v", name, p),
			http.StatusInternalServerError)
		return
	}

	responseErr := s.response(w, http.StatusOK, p)
	if responseErr != nil {
		log.Printf("failed to write http response %v", responseErr)
	}
}

func (s *Service) response(w http.ResponseWriter, statusCode int, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal response JSON: %w", err)
	}

	w.WriteHeader(statusCode)

	if _, writeErr := w.Write(jsonData); writeErr != nil {
		return fmt.Errorf("failed to write response JSON: %w", writeErr)
	}

	return nil
}

func (s *Service) renderError(writer http.ResponseWriter, errorMessage string, statusCode int) {
	errResponse := &ErrorResponse{
		Message: errorMessage,
	}
	errJSON, err := json.Marshal(errResponse)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(writer, string(errJSON), statusCode)
}
