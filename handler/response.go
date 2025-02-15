package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

func sendResponse(w http.ResponseWriter, response APIResponse) error {

	for k, v := range response.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(response.StatusCode)

	if response.Body != nil {
		body, err := json.Marshal(response.Body)
		w.Write(body)
		if err != nil {
			return fmt.Errorf("failed to send response %w", err)
		}
	}
	return nil
}
