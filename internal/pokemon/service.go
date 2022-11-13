package pokemon

import (
	"context"
	"net/http"
	"pokedex-clone/internal/api"
	"pokedex-clone/internal/storage"

	"github.com/gin-gonic/gin"
)

type Service struct {
	StorageAPI    *storage.Store
	PokeAPIClient *api.Client
}

func NewService(storage *storage.Store, pokeAPIClient *api.Client) *Service {
	return &Service{
		StorageAPI:    storage,
		PokeAPIClient: pokeAPIClient,
	}
}

func (s *Service) Get(c *gin.Context) {
	var req PokemonNameUri
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := req.Name

	pokemon, err := s.PokeAPIClient.GetSpecies(context.Background(), name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
	}
	// if !s.StorageAPI.Exist(name) {
	// 	c.JSON(http.StatusNotFound, fmt.Sprintf("missing pokemon with name %v", name))
	// 	return
	// }

	// pokemon, ok := s.StorageAPI.Load(name)
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, fmt.Sprintf("failed to get pokemon %v", name))
	// 	return
	// }

	// p, ok := pokemon.(api.PokemonSpecies)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, fmt.Sprintf("unexpected object returned for %v, %v", name, p))
	// 	return
	// }

	c.JSON(http.StatusOK, pokemon)
}

// func (s *Service) GetTranslated(c *gin.Context) {
// 	var req PokemonNameUri
// 	if err := c.ShouldBindUri(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	name := req.Name

// 	// todo remove
// 	err := s.StorageAPI.Save("mewtwo", Pokemon{
// 		Name:        "mewtwo",
// 		Description: "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments.",
// 		Habitat:     "rare",
// 		IsLegendary: true,
// 	})
// 	if err != nil {
// 		log.Println("oley")
// 	}

// 	if !s.StorageAPI.Exist(name) {
// 		c.JSON(http.StatusNotFound, fmt.Sprintf("missing pokemon with name %v", name))
// 		return
// 	}

// 	pokemon, ok := s.StorageAPI.Load(name)
// 	if !ok {
// 		c.JSON(http.StatusBadRequest, fmt.Sprintf("failed to get pokemon %v", name))
// 		return
// 	}

// 	p, ok := pokemon.(Pokemon)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, fmt.Sprintf("unexpected object returned for %v, %v", name, p))
// 		return
// 	}

// 	c.JSON(http.StatusOK, p)
// }

// func (s *Service) response(w http.ResponseWriter, statusCode int, data interface{}) error {
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal response JSON: %w", err)
// 	}

// 	w.WriteHeader(statusCode)

// 	if _, writeErr := w.Write(jsonData); writeErr != nil {
// 		return fmt.Errorf("failed to write response JSON: %w", writeErr)
// 	}

// 	return nil
// }

// func (s *Service) renderError(writer http.ResponseWriter, errorMessage string, statusCode int) {
// 	errResponse := &ErrorResponse{
// 		Message: errorMessage,
// 	}
// 	errJSON, err := json.Marshal(errResponse)
// 	if err != nil {
// 		http.Error(writer, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Error(writer, string(errJSON), statusCode)
// }
