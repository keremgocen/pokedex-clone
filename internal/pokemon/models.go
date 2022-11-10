package pokemon

// "name": "mewtwo",
// "description": "It was created by a scientist after years of horrific gene
// splicing and DNA engineering experiments.",
// "habitat": "rare",
// "isLegendary": true

type Pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary bool   `json:"is_legendary"`
}

// type PokedexData struct {
// 	Data Pokemon `json:"data"`
// }
