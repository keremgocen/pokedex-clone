package pokemon

type ErrorResponse struct {
	Message string `json:"message"`
}

type PokemonNameUri struct {
	Name string `uri:"name" binding:"required"`
}
