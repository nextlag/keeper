package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// AddLogin adds a new login for the current user.
func (c *Controller) AddLogin(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("getUserFromCtx", l.ErrAttr(err))
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

// GetLogins - get the logins of the current user
func (c *Controller) GetLogins(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("getUserFromCtx", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
	}

	userLogins, err := c.uc.GetLogins(r.Context(), currentUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(userLogins) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(userLogins); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DelLogin handles the deletion of a login by its UUID.
func (c *Controller) DelLogin(w http.ResponseWriter, r *http.Request) {
	loginUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "uuid", loginUUID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	if err = c.uc.DelLogin(r.Context(), loginUUID, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// UpdateLogin handles the update of a login by its UUID.
func (c *Controller) UpdateLogin(w http.ResponseWriter, r *http.Request) {
	loginUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "uuid", loginUUID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var payloadLogin entity.Login
	if err = json.NewDecoder(r.Body).Decode(&payloadLogin); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payloadLogin.ID = loginUUID

	if err = c.uc.UpdateLogin(r.Context(), &payloadLogin, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
