package v1

type HTTPError struct {
	Message string `json:"message" example:"error description"`
}
