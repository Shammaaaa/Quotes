package handlers

import (
	"Quotes/internal/models/api"
	"encoding/json"
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, message, code string) {
	log.Printf("Error: %s (status: %d, code: %s)", message, status, code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(api.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
