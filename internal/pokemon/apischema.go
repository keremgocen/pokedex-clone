package pokemon

type ErrorResponse struct {
	Message string `json:"message"`
}

type NameURI struct {
	Name string `uri:"name" binding:"required"`
}
