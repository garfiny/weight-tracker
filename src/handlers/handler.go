package handlers

import (
	"encoding/json"
	"net/http"
)

// HandleRequest processes incoming requests and returns responses.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{"message": "Hello, World!"}
	json.NewEncoder(w).Encode(response)
}

// Additional helper functions can be defined here for handling specific routes.
