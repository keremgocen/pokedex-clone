package pokemon

import (
	"context"
	"log"
	"net/http"
	"pokedex-clone/pkg/api"
	"pokedex-clone/pkg/storage"

	"github.com/gin-gonic/gin"
)

const (
	ISO639ENGString = "en"
)

type Service struct {
	StorageAPI      *storage.Store
	PokeAPI         api.PokeAPI
	TranslationsAPI api.TranslationsAPI
}

func NewService(storage *storage.Store, pokeAPI api.PokeAPI, translationsAPI api.TranslationsAPI) *Service {
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

	if cachedPokemon, ok := s.StorageAPI.Load(name); ok {
		c.JSON(http.StatusOK, cachedPokemon)
		return
	}

	pokemonSpecies, err := s.PokeAPI.GetSpecies(context.Background(), name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
	}

	descriptionText, _ := getFirstEnglishFlavorText(pokemonSpecies.FlaworTextEntries)

	pokemon := Pokemon{
		Description: descriptionText,
		IsLegendary: pokemonSpecies.IsLegendary,
		Habitat:     pokemonSpecies.Habitat.Name,
		Name:        pokemonSpecies.Name,
	}

	if cacheErr := s.StorageAPI.Save(name, &pokemon); cacheErr != nil {
		log.Printf("failed to save %s in cache: [%v]", name, cacheErr.Error())
	}

	c.JSON(http.StatusOK, pokemon)
}

func (s *Service) GetTranslated(c *gin.Context) {
	var req NameURI
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := req.Name

	// we can potentially avoid this API call if Get was called before
	pokemonSpec, err := s.PokeAPI.GetSpecies(context.Background(), name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
	}

	// check description text and maybe skip API calls
	descriptionText, languageCode := getFirstEnglishFlavorText(pokemonSpec.FlaworTextEntries)
	if languageCode != ISO639ENGString {
		c.JSON(http.StatusOK, Pokemon{
			Description: descriptionText,
			IsLegendary: pokemonSpec.IsLegendary,
			Habitat:     pokemonSpec.Habitat.Name,
			Name:        pokemonSpec.Name,
		})
		return
	}

	var translationType api.TranslationType
	if pokemonSpec.Habitat.Name == "cave" || pokemonSpec.IsLegendary {
		// use yoda translation
		translationType = api.TTypeYoda
	} else {
		// use Shakespeare translation
		translationType = api.TTypeShakespeare
	}

	if cachedPokemonWithTrans, ok := s.StorageAPI.Load(name + string(translationType)); ok {
		c.JSON(http.StatusOK, cachedPokemonWithTrans)
		return
	}

	response, tErr := s.TranslationsAPI.GetTranslation(context.Background(), name, descriptionText, translationType)
	if tErr == nil && response.Success.Total > 0 {
		descriptionText = response.Contents.Translated
	}

	p := Pokemon{
		Description: descriptionText,
		IsLegendary: pokemonSpec.IsLegendary,
		Habitat:     pokemonSpec.Habitat.Name,
		Name:        pokemonSpec.Name,
	}

	if cacheErr := s.StorageAPI.Save(name+string(translationType), &p); cacheErr != nil {
		log.Printf("failed to save %s in cache: [%v]", name, cacheErr.Error())
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
