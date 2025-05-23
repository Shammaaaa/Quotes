package api

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	ErrCodeInvalidRequest = "invalid_request"
	ErrCodeNotFound       = "not_found"
	ErrCodeInternal       = "internal_server"
)
