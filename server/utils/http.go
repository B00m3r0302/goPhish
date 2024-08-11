package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse represents a standard JSON response structure
type JSONResponse struct {
	Status	int			'json:"status"'
	Message	string		'json:"message"'
	Data	interface[]	'json."data,omitempty"'
}

// SendJSONResponse sends a JSON response with the given status, message and data
func SendJSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	response := JSONResponse{
		Status:		status,
		Message:	message,
		Data:		data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// HandleCORS sets the CORS headers for the response
func HandleCORS(w http.ResponseWriter, allowedOrigins []string) {
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0]) // Simplified for a single origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// HandleOptionsRequest handles preflight OPTIONS requests for CORS
func HandleOptionsRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}