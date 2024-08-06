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

// AddLogin godoc
// @Summary Add a new login
// @Description Upload a new login for the current user
// @Tags logins
// @Accept json
// @Produce json
// @Param login body entity.Login true "Login data"
// @Success 202 {object} entity.Login
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/logins [post]
func (c *Controller) AddLogin(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	var payloadLogin entity.Login

	if err = json.NewDecoder(r.Body).Decode(&payloadLogin); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	if err = c.uc.AddLogin(r.Context(), &payloadLogin, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(payloadLogin); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// GetLogins godoc
// @Summary Get all logins for the current user
// @Description Retrieve all logins for the current user
// @Tags logins
// @Produce json
// @Success 200 {array} entity.Login
// @Success 204 "No content"
// @Failure 500 {object} response
// @Router /user/logins [get]
func (c *Controller) GetLogins(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
	}

	userLogins, err := c.uc.GetLogins(r.Context(), currentUser)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
	}

	if len(userLogins) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(userLogins); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// UpdateLogin godoc
// @Summary Update a login by UUID
// @Description Update a specific login identified by its UUID
// @Tags logins
// @Accept json
// @Produce json
// @Param id path string true "Login UUID"
// @Param login body entity.Login true "Updated login data"
// @Success 202 {string} string "Update accepted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/logins/{id} [patch]
func (c *Controller) UpdateLogin(w http.ResponseWriter, r *http.Request) {
	loginUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "uuid", loginUUID)
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	var payloadLogin entity.Login
	if err = json.NewDecoder(r.Body).Decode(&payloadLogin); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}
	payloadLogin.ID = loginUUID

	if err = c.uc.UpdateLogin(r.Context(), &payloadLogin, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte(jsonResponse("update accepted"))); err != nil {
		return
	}
}

// DelLogin godoc
// @Summary Delete a login by UUID
// @Description Delete a specific login identified by its UUID
// @Tags logins
// @Param id path string true "Login UUID"
// @Success 202 {string} string "Delete accepted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/logins/{id} [delete]
func (c *Controller) DelLogin(w http.ResponseWriter, r *http.Request) {
	loginUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "uuid", loginUUID)
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	if err = c.uc.DelLogin(r.Context(), loginUUID, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte(jsonResponse("delete accepted"))); err != nil {
		return
	}
}
