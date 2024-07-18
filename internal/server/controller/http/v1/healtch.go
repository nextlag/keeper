package v1

import (
	"net/http"
)

// HealthCheck checks the health of the application.
func (c *Controller) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	err := c.uc.HealthCheck()
	if err != nil {
		http.Error(w, "Application not available", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte(jsonResponse("connected"))); err != nil {
		return
	}
}
