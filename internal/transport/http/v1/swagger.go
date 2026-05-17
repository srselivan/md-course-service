package v1

// ErrorResponse describes an API error message returned by Fiber handlers.
type ErrorResponse struct {
	Message string `json:"message" example:"bad request"`
}
