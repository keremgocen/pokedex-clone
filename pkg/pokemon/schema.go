package pokemon

type NameURI struct {
	Name string `uri:"name" binding:"required,alpha"`
}

type Pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary bool   `json:"is_legendary"`
}
