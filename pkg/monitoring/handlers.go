package monitoring

import (
	"encoding/json"
	"net/http"
)

type HealthStatus struct {
	Status string `json:"status"`
}

// Check if service is healthy
func GetHealth(w http.ResponseWriter, r *http.Request) {
	data := HealthStatus{}
	data.Status = "UP"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// Check if service is healthy
func GetInfo(w http.ResponseWriter, r *http.Request) {
	data := HealthStatus{}
	data.Status = "UP"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
