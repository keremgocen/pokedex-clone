package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pokedex-clone/internal/api"
	"pokedex-clone/internal/pokemon"
	"pokedex-clone/internal/storage"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	ctxTimeout    = 5 * time.Second
	serverTimeout = 3 * time.Second
	pokeAPIURL    = "https://pokeapi.co/api/v2/pokemon-species/"
)

func run() (*gin.Engine, *pokemon.Service) {
	storageAPI := storage.NewStore()
	pokeAPI := api.NewClient(pokeAPIURL, serverTimeout, storageAPI)
	service := pokemon.NewService(storageAPI, pokeAPI)

	router := gin.Default()
	router.GET("/pokemon/:name", service.Get)
	// router.GET("/pokemon/translated/:name", service.GetTranslated)

	return router, service
}

func TestGetPokemon(t *testing.T) {
	router, _ := run()

	defaultPokemon := &api.PokemonSpecies{
		Name: "mewtwo",
		FlaworTextEntries: []api.FlaworText{
			{
				FlaworText: "some text here",
				Language: api.NamedAPIResource{
					Name: "en",
					URL:  "",
				},
				Version: api.NamedAPIResource{
					Name: "1",
					URL:  "",
				},
			},
		},
	}

	tests := map[string]struct {
		name          string
		wantStatus    int
		wantErr       error
		cachedPokemon *api.PokemonSpecies
	}{
		"get pokemon returns 200 with expected name": {
			name:          "mewtwo",
			wantStatus:    http.StatusOK,
			wantErr:       nil,
			cachedPokemon: defaultPokemon,
		},
		// "missing pokemon returns 404": {
		// 	name:          "unknown",
		// 	wantStatus:    http.StatusNotFound,
		// 	wantErr:       nil,
		// 	cachedPokemon: nil,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/pokemon/"+tc.name, nil)
			assert.Equal(t, err, tc.wantErr)

			// storageErr := service.StorageAPI.Save(tc.name, tc.cachedPokemon)
			// assert.Nil(t, storageErr)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.wantStatus, rr.Code)

			var pokemon api.PokemonSpecies
			err = json.Unmarshal(rr.Body.Bytes(), &pokemon)
			assert.Nil(t, err)

			assert.Equal(t, tc.name, pokemon.Name)
		})
	}
}

// func TestGetUnexpectedMethodReturns404(t *testing.T) {
// 	router, _ := run()

// 	req, err := http.NewRequest(http.MethodPatch, "/pokemon/mewtwo", nil)
// 	assert.Nil(t, err)

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusNotFound, rr.Code)
// }

// func TestGetMalformedPokemonReturns500(t *testing.T) {
// 	router, service := run()

// 	req, err := http.NewRequest(http.MethodGet, "/pokemon/malformedPokemon", nil)
// 	assert.Nil(t, err)

// 	saveErr := service.StorageAPI.Save("malformedPokemon", "malformedPokemonObjectStr")
// 	assert.Nil(t, saveErr)

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
// }
