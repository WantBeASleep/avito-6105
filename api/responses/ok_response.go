package responses

import (
	"encoding/json"
	"net/http"
)

func OkJSON(w http.ResponseWriter, statusCode int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(resp)
}
