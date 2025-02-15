package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Errors string `json:"errors,omitempty"`
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) error {

	body, err := json.Marshal(ErrorResponse{
		Errors: errorMessage,
	})
	if err != nil {
		http.Error(w, errorMessage, statusCode)
		return fmt.Errorf("failed to send error response %w", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	w.Write(body)
	return nil
}
