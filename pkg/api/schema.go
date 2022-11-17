package api

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PokemonSpecies represents the returned payload from pokeapi.
type PokemonSpecies struct {
	Name              string           `json:"name"`
	FlavorTextEntries []FlavorText     `json:"flavor_text_entries"`
	Habitat           NamedAPIResource `json:"habitat"`
	IsLegendary       bool             `json:"is_legendary"`
}

type FlavorText struct {
	FlavorText string           `json:"flavor_text"`
	Language   NamedAPIResource `json:"language"`
	Version    NamedAPIResource `json:"version"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type TranslationText struct {
	Text string `json:"text"`
}

type Success struct {
	Total int `json:"total"`
}

type Contents struct {
	Translated  string `json:"translated"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

type TranslateAPIResponse struct {
	Success  Success  `json:"success"`
	Contents Contents `json:"contents"`
}
