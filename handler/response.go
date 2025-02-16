package handler

import (
	"encoding/json"
	"fmt"
	"merch-store/models"
	"net/http"
)

func sendResponse(w http.ResponseWriter, response models.APIResponse) error {

	for k, v := range response.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(response.StatusCode)

	if response.Body != nil {
		body, err := json.Marshal(response.Body)
		w.Write(body)
		if err != nil {
			return fmt.Errorf("failed to marshal response body: %w", err)
		}
	}

	return nil
}
