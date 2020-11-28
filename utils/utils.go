package utils

import (
	"encoding/json"
	"net/http"
)

func WriteMessageResponse(w *http.ResponseWriter, statusCode int, message string) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(statusCode)

	respMsg := make(map[string]string)
	respMsg["message"] = message
	json.NewEncoder(*w).Encode(respMsg)
}
