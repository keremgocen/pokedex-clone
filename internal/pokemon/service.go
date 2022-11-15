package pokemon

import (
	"context"
	"net/http"
	"pokedex-clone/internal/api"
	"pokedex-clone/internal/storage"

	"github.com/gin-gonic/gin"
)

const (
	ISO639ENGString = "en"
)

type Service struct {
	StorageAPI      *storage.Store
	PokeAPI         api.Poke
	TranslationsAPI api.Translations
}

func NewService(storage *storage.Store, pokeAPI api.Poke, translationsAPI api.Translations) *Service {
	return &Service{
		StorageAPI:      storage,
		PokeAPI:         pokeAPI,
		TranslationsAPI: translationsAPI,
	}
}

func (s *Service) Get(c *gin.Context) {
	var req NameURI
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := req.Name

	pokemon, err := s.PokeAPI.GetSpecies(context.Background(), name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
	}

	descriptionText, _ := getFirstEnglishFlavorText(pokemon.FlaworTextEntries)
	c.JSON(http.StatusOK, Pokemon{
		Description: descriptionText,
		IsLegendary: pokemon.IsLegendary,
		Habitat:     pokemon.Habitat.Name,
		Name:        pokemon.Name,
	})
}

func (s *Service) GetTranslated(c *gin.Context) {
	var req NameURI
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := req.Name

	pokemon, err := s.PokeAPI.GetSpecies(context.Background(), name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
	}

	// check description text and maybe skip API calls
	descriptionText, languageCode := getFirstEnglishFlavorText(pokemon.FlaworTextEntries)
	if languageCode != ISO639ENGString {
		c.JSON(http.StatusOK, Pokemon{
			Description: descriptionText,
			IsLegendary: pokemon.IsLegendary,
			Habitat:     pokemon.Habitat.Name,
			Name:        pokemon.Name,
		})
		return
	}

	if pokemon.Habitat.Name == "cave" || pokemon.IsLegendary {
		// use yoda translation
		response, tErr := s.TranslationsAPI.GetTranslation(context.Background(), name, descriptionText, api.TTypeYoda)
		if tErr == nil && response.Success.Total > 0 {
			descriptionText = response.Contents.Translated
		}
	} else {
		// use Shakespeare translation
		response, tErr := s.TranslationsAPI.GetTranslation(context.Background(), name, descriptionText, api.TTypeShakespeare)
		if tErr == nil && response.Success.Total > 0 {
			descriptionText = response.Contents.Translated
		}
	}

	p := Pokemon{
		Description: descriptionText,
		IsLegendary: pokemon.IsLegendary,
		Habitat:     pokemon.Habitat.Name,
		Name:        pokemon.Name,
	}

	c.JSON(http.StatusOK, p)
}

func getFirstEnglishFlavorText(entries []api.FlaworText) (string, string) {
	if len(entries) > 0 {
		for _, entry := range entries {
			if entry.Language.Name == ISO639ENGString {
				return entry.FlaworText, ISO639ENGString
			}
		}

		return entries[0].FlaworText, entries[0].Language.Name
	}

	return "", ""
}
