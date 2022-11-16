package pokemon_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pokedex-clone/internal/api"
	"pokedex-clone/internal/api/mocks"
	"pokedex-clone/internal/pokemon"
	"pokedex-clone/internal/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetPokemon(t *testing.T) {
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
		name              string
		wantStatus        int
		wantErr           error
		getSpeciesReturns *api.PokemonSpecies
		getSpeciesErr     error
	}{
		"get pokemon returns 200 with expected name": {
			name:              "mewtwo",
			wantStatus:        http.StatusOK,
			wantErr:           nil,
			getSpeciesReturns: defaultPokemon,
			getSpeciesErr:     nil,
		},
		"missing pokemon returns 404 when get species can't find it": {
			name:              "missing",
			wantStatus:        http.StatusNotFound,
			wantErr:           nil,
			getSpeciesReturns: nil,
			getSpeciesErr:     gin.Error{Err: fmt.Errorf("not found")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storageAPI := storage.NewStore()

			ctrl := gomock.NewController(t)

			mockPokeAPI := mocks.NewMockPokeAPI(ctrl)
			mockTranslationsAPI := mocks.NewMockTranslationsAPI(ctrl)
			service := pokemon.NewService(storageAPI, mockPokeAPI, mockTranslationsAPI)

			router := gin.Default()
			router.GET("/pokemon/:name", service.Get)

			req, err := http.NewRequest(http.MethodGet, "/pokemon/"+tc.name, nil)
			assert.Equal(t, err, tc.wantErr)

			mockPokeAPI.EXPECT().GetSpecies(gomock.Any(), tc.name).Return(tc.getSpeciesReturns, tc.getSpeciesErr)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.wantStatus, rr.Code)

			var pokemon pokemon.Pokemon
			err = json.Unmarshal(rr.Body.Bytes(), &pokemon)
			assert.Nil(t, err)
		})
	}
}

func TestGetPokemonFails(t *testing.T) {
	tests := map[string]struct {
		name       string
		wantStatus int
		wantErr    error
	}{
		"get pokemon fails 400 with non-alphanumeric name": {
			name:       "123",
			wantStatus: http.StatusBadRequest,
			wantErr:    nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storageAPI := storage.NewStore()

			ctrl := gomock.NewController(t)

			mockPokeAPI := mocks.NewMockPokeAPI(ctrl)
			mockTranslationsAPI := mocks.NewMockTranslationsAPI(ctrl)
			service := pokemon.NewService(storageAPI, mockPokeAPI, mockTranslationsAPI)

			router := gin.Default()
			router.GET("/pokemon/:name", service.Get)

			req, err := http.NewRequest(http.MethodGet, "/pokemon/"+tc.name, nil)
			assert.Equal(t, err, tc.wantErr)

			mockPokeAPI.EXPECT().GetSpecies(gomock.Any(), tc.name).Times(0)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.wantStatus, rr.Code)
		})
	}
}

func TestPokemonTranslation200(t *testing.T) {
	yodaTranslatedPokemon := pokemon.Pokemon{
		Name:        "mewtwo",
		Description: "Created by a scientist after years of horrific gene splicing and dna engineering experiments, it was.",
		Habitat:     "cave",
		IsLegendary: false,
	}

	cavePokemonSpecies := &api.PokemonSpecies{
		Name: "mewtwo",
		FlaworTextEntries: []api.FlaworText{
			{
				FlaworText: "It was created by a scientist after years of horrific gene splicing and dna engineering experiments.",
				Language: api.NamedAPIResource{
					Name: "en",
					URL:  "",
				},
			},
		},
		Habitat: api.NamedAPIResource{
			Name: "cave",
			URL:  "",
		},
		IsLegendary: false,
	}

	yodaTranslatedLegendaryPokemon := pokemon.Pokemon{
		Name:        "mewlegend",
		Description: "Created by a scientist after years of horrific gene splicing and dna engineering experiments, it was.",
		Habitat:     "indoors",
		IsLegendary: true,
	}

	legendaryPokemonSpecies := &api.PokemonSpecies{
		Name: "mewlegend",
		FlaworTextEntries: []api.FlaworText{
			{
				FlaworText: "It was created by a scientist after years of horrific gene splicing and dna engineering experiments.",
				Language: api.NamedAPIResource{
					Name: "en",
					URL:  "",
				},
			},
		},
		Habitat: api.NamedAPIResource{
			Name: "indoors",
			URL:  "",
		},
		IsLegendary: true,
	}

	shakespeareTranslatedPokemon := pokemon.Pokemon{
		Name:        "thepoet",
		Description: "Some text.",
		Habitat:     "somewhere",
		IsLegendary: false,
	}

	shakespeareanPokemonSpecies := &api.PokemonSpecies{
		Name: "thepoet",
		FlaworTextEntries: []api.FlaworText{
			{
				FlaworText: "Shakespeare translation of some text",
				Language: api.NamedAPIResource{
					Name: "en",
					URL:  "",
				},
			},
		},
		Habitat: api.NamedAPIResource{
			Name: "somewhere",
			URL:  "",
		},
		IsLegendary: false,
	}

	tests := map[string]struct {
		name                  string
		wantStatus            int
		wantErr               error
		getSpeciesReturns     *api.PokemonSpecies
		getSpeciesErr         error
		getTranslationReturns *api.TranslateAPIResponse
		getTranslationErr     error
		translatedPokemon     *pokemon.Pokemon
		translationType       api.TranslationType
	}{
		"pokemon in cave, not legendary still gets yoda translation 200": {
			name:              "mewtwo",
			wantStatus:        http.StatusOK,
			wantErr:           nil,
			getSpeciesReturns: cavePokemonSpecies,
			getSpeciesErr:     nil,
			getTranslationReturns: &api.TranslateAPIResponse{
				Success: api.Success{
					Total: 1,
				},
				Contents: api.Contents{
					Translated: yodaTranslatedPokemon.Description,
				},
			},
			getTranslationErr: nil,
			translatedPokemon: &yodaTranslatedPokemon,
			translationType:   api.TTypeYoda,
		},
		"legendary pokemon also gets yoda translation 200": {
			name:              "mewlegend",
			wantStatus:        http.StatusOK,
			wantErr:           nil,
			getSpeciesReturns: legendaryPokemonSpecies,
			getSpeciesErr:     nil,
			getTranslationReturns: &api.TranslateAPIResponse{
				Success: api.Success{
					Total: 1,
				},
				Contents: api.Contents{
					Translated: yodaTranslatedLegendaryPokemon.Description,
				},
			},
			getTranslationErr: nil,
			translatedPokemon: &yodaTranslatedLegendaryPokemon,
			translationType:   api.TTypeYoda,
		},
		"non-legendary, non-cave pokemon gets shakespeare translation 200": {
			name:              "mewlegend",
			wantStatus:        http.StatusOK,
			wantErr:           nil,
			getSpeciesReturns: shakespeareanPokemonSpecies,
			getSpeciesErr:     nil,
			getTranslationReturns: &api.TranslateAPIResponse{
				Success: api.Success{
					Total: 1,
				},
				Contents: api.Contents{
					Translated: shakespeareTranslatedPokemon.Description,
				},
			},
			getTranslationErr: nil,
			translatedPokemon: &shakespeareTranslatedPokemon,
			translationType:   api.TTypeShakespeare,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storageAPI := storage.NewStore()

			ctrl := gomock.NewController(t)

			mockPokeAPI := mocks.NewMockPokeAPI(ctrl)
			mockTranslationsAPI := mocks.NewMockTranslationsAPI(ctrl)
			service := pokemon.NewService(storageAPI, mockPokeAPI, mockTranslationsAPI)

			router := gin.Default()
			router.GET("/pokemon/translated/:name", service.GetTranslated)

			req, err := http.NewRequest(http.MethodGet, "/pokemon/translated/"+tc.name, nil)
			assert.Equal(t, err, tc.wantErr)

			mockPokeAPI.EXPECT().GetSpecies(gomock.Any(), tc.name).Return(tc.getSpeciesReturns, tc.getSpeciesErr)
			mockTranslationsAPI.EXPECT().GetTranslation(
				gomock.Any(),
				tc.name,
				tc.getSpeciesReturns.FlaworTextEntries[0].FlaworText,
				tc.translationType).Return(tc.getTranslationReturns, tc.getTranslationErr)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.wantStatus, rr.Code)

			var pokemon pokemon.Pokemon
			err = json.Unmarshal(rr.Body.Bytes(), &pokemon)
			assert.Nil(t, err)

			assert.Equal(t, tc.translatedPokemon.Name, pokemon.Name)
			assert.Equal(t, tc.translatedPokemon.Description, pokemon.Description)
			assert.Equal(t, tc.translatedPokemon.Habitat, pokemon.Habitat)
			assert.Equal(t, tc.translatedPokemon.IsLegendary, pokemon.IsLegendary)
		})
	}
}

func TestFailedPokemonTranslationReturnsDefaults(t *testing.T) {
	defaultPokemon := pokemon.Pokemon{
		Description: "Translated description.",
	}

	frenchPokemonSpecies := &api.PokemonSpecies{
		Name: "mewtwo",
		FlaworTextEntries: []api.FlaworText{
			{
				FlaworText: "Some really French text.",
				Language: api.NamedAPIResource{
					Name: "fr",
					URL:  "",
				},
			},
		},
		Habitat: api.NamedAPIResource{
			Name: "cave",
			URL:  "",
		},
		IsLegendary: false,
	}

	engPokemonSpecies := &api.PokemonSpecies{
		Name: "mewtwo",
		FlaworTextEntries: []api.FlaworText{
			{
				FlaworText: "Some English text.",
				Language: api.NamedAPIResource{
					Name: "en",
					URL:  "",
				},
			},
		},
		Habitat: api.NamedAPIResource{
			Name: "cave",
			URL:  "",
		},
		IsLegendary: false,
	}

	unknownPokemonSpecies := &api.PokemonSpecies{
		Name:              "mewtwo",
		FlaworTextEntries: []api.FlaworText{{}},
		Habitat: api.NamedAPIResource{
			Name: "cave",
			URL:  "",
		},
		IsLegendary: false,
	}

	tests := map[string]struct {
		name                  string
		wantStatus            int
		wantErr               error
		getSpeciesReturns     *api.PokemonSpecies
		getSpeciesErr         error
		getTranslationReturns *api.TranslateAPIResponse
		getTranslationErr     error
		translationType       api.TranslationType
		translationCallCount  int
		expectedTranslation   string
	}{
		"description is not translated when English description is missing": {
			name:              "mewtwo",
			wantStatus:        http.StatusOK,
			wantErr:           nil,
			getSpeciesReturns: frenchPokemonSpecies,
			getSpeciesErr:     nil,
			getTranslationReturns: &api.TranslateAPIResponse{
				Success: api.Success{
					Total: 1,
				},
				Contents: api.Contents{
					Translated: defaultPokemon.Description,
				},
			},
			getTranslationErr:    nil,
			translationType:      api.TTypeYoda,
			translationCallCount: 0,
			expectedTranslation:  frenchPokemonSpecies.FlaworTextEntries[0].FlaworText,
		},
		"pokemon species value is is used when translateAPI returns err": {
			name:                  "mewtwo",
			wantStatus:            http.StatusOK,
			wantErr:               nil,
			getSpeciesReturns:     engPokemonSpecies,
			getSpeciesErr:         nil,
			getTranslationReturns: nil,
			getTranslationErr:     gin.Error{},
			translationType:       api.TTypeYoda,
			translationCallCount:  1,
			expectedTranslation:   engPokemonSpecies.FlaworTextEntries[0].FlaworText,
		},
		"pokemon with no description entries": {
			name:              "mewtwo",
			wantStatus:        http.StatusOK,
			wantErr:           nil,
			getSpeciesReturns: unknownPokemonSpecies,
			getSpeciesErr:     nil,
			getTranslationReturns: &api.TranslateAPIResponse{
				Success: api.Success{
					Total: 0,
				},
			},
			getTranslationErr:    nil,
			translationType:      api.TTypeYoda,
			translationCallCount: 0,
			expectedTranslation:  "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storageAPI := storage.NewStore()

			ctrl := gomock.NewController(t)

			mockPokeAPI := mocks.NewMockPokeAPI(ctrl)
			mockTranslationsAPI := mocks.NewMockTranslationsAPI(ctrl)
			service := pokemon.NewService(storageAPI, mockPokeAPI, mockTranslationsAPI)

			router := gin.Default()
			router.GET("/pokemon/translated/:name", service.GetTranslated)

			req, err := http.NewRequest(http.MethodGet, "/pokemon/translated/"+tc.name, nil)
			assert.Equal(t, err, tc.wantErr)

			mockPokeAPI.EXPECT().GetSpecies(gomock.Any(), tc.name).Return(tc.getSpeciesReturns, tc.getSpeciesErr)
			mockTranslationsAPI.EXPECT().GetTranslation(
				gomock.Any(),
				tc.name,
				tc.getSpeciesReturns.FlaworTextEntries[0].FlaworText,
				tc.translationType).Times(tc.translationCallCount).Return(tc.getTranslationReturns, tc.getTranslationErr)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.wantStatus, rr.Code)

			var pokemon pokemon.Pokemon
			err = json.Unmarshal(rr.Body.Bytes(), &pokemon)
			assert.Nil(t, err)

			assert.Equal(t, tc.expectedTranslation, pokemon.Description)
		})
	}
}
