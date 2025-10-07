// Package handlers contains the HTTP request handlers for the API.
package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheck is the handler for the GET /healthz endpoint.
// It provides a basic health check for the service.
func HealthCheck(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}