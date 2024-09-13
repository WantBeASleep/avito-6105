package responses

import (
	"encoding/json"
	"net/http"
)

type ErrAnswer struct {
	Reason string `json:"reason"`
}

func Error(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func ErrorJSON(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrAnswer{Reason: err.Error()}
	json.NewEncoder(w).Encode(errorResponse)
}
