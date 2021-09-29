package http

import (
	"net/http"

	"github.com/julb/go/pkg/monitoring"
)

// HTTP handler for health endpoint to return the system status as JSON format.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		contributor := monitoring.GetHealthContributor().ToMap()
		healthStatus := contributor["status"].(monitoring.SystemStatus)
		if healthStatus == monitoring.Up || healthStatus == monitoring.Partial || healthStatus == monitoring.Unknown {
			WriteJSON(r, w, http.StatusOK, contributor)
		} else {
			WriteJSON(r, w, http.StatusServiceUnavailable, contributor)
		}
	default:
		WriteNotFound(r, w)
	}
}

// HTTP handler for info endpoint to return the system general information as JSON format.
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		contributor := monitoring.GetInfoContributor().ToMap()
		WriteJSON(r, w, http.StatusOK, contributor)
	default:
		WriteNotFound(r, w)
	}
}

// HTTP handler for info endpoint to return the system general information as JSON format.
func RuntimeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		contributor := monitoring.GetRuntimeContributor().ToMap()
		WriteJSON(r, w, http.StatusOK, contributor)
	default:
		WriteNotFound(r, w)
	}
}
