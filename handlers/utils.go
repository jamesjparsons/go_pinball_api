package handlers

import (
	"encoding/json"
	"net/http"
)

// JSONResponse represents a standard JSON response structure
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SendJSON sends a JSON response
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// SendError sends an error JSON response
func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, JSONResponse{
		Success: false,
		Message: message,
	})
}

// SendSuccess sends a success JSON response
func SendSuccess(w http.ResponseWriter, data interface{}) {
	SendJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    data,
	})
} 