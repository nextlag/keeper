package v1

import (
	"encoding/json"
	"net/http"

	"github.com/nextlag/keeper/internal/entity"
)

// AddLogin adds a new login for the current user.
func (c *Controller) AddLogin(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var payloadLogin entity.Login

	if err = json.NewDecoder(r.Body).Decode(&payloadLogin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.uc.AddLogin(r.Context(), &payloadLogin, currentUser.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(payloadLogin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
