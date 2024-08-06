package v1

import (
	"net/http"

	"github.com/nextlag/keeper/pkg/logger/l"
)

// HealthCheck godoc
// @Summary Check the health of the application
// @Description Endpoint to check if the application is running correctly
// @Tags health
// @Produce json
// @Success 200 {string} string "connected"
// @Failure 500 {object} response
// @Router /ping [get]
func (c *Controller) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	err := c.uc.HealthCheck()
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte(jsonResponse("connected"))); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		return
	}
}
