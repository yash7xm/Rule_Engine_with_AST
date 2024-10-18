package routes

import (
	"encoding/json"
	"net/http"
)

// Helper function to send detailed error responses
func SendErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := map[string]interface{}{
		"message": message,
		"error":   err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Helper function to send detailed success responses
func SendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data map[string]interface{}) {
	response := map[string]interface{}{
		"message": message,
		"data":    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
