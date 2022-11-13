package pokemon

type Pokemon struct {
	Name              string `json:"name"`
	FlaworTextEntries string `json:"flavor_text_entries"`
	Habitat           string `json:"habitat"`
	IsLegendary       bool   `json:"is_legendary"`
}
