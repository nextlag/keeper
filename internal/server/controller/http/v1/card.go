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

// AddCard godoc
// @Summary Add a new card
// @Description Upload a new card for the current user
// @Tags cards
// @Accept json
// @Produce json
// @Param card body entity.Card true "Card data"
// @Success 202 {object} entity.Card
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/cards [post]
func (c *Controller) AddCard(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
	}

	var payloadCard *entity.Card
	if err = json.NewDecoder(r.Body).Decode(&payloadCard); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	if err = c.uc.AddCard(r.Context(), payloadCard, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(payloadCard); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// GetCards godoc
// @Summary Get all cards for the current user
// @Description Retrieve all cards for the current user
// @Tags cards
// @Produce json
// @Success 200 {array} entity.Card
// @Success 204 "No content"
// @Failure 500 {object} response
// @Router /user/cards [get]
func (c *Controller) GetCards(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
	}

	userCards, err := c.uc.GetCards(r.Context(), currentUser)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
	}

	if len(userCards) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(userCards); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// UpdateCard godoc
// @Summary Update a card by UUID
// @Description Update a specific card identified by its UUID
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card UUID"
// @Param card body entity.Card true "Updated card data"
// @Success 202 {string} string "Update accepted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/cards/{id} [patch]
func (c *Controller) UpdateCard(w http.ResponseWriter, r *http.Request) {
	cardUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	var payloadCard *entity.Card
	if err = json.NewDecoder(r.Body).Decode(&payloadCard); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	payloadCard.ID = cardUUID
	if err = c.uc.UpdateCard(r.Context(), payloadCard, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte("Update accepted")); err != nil {
		return
	}
}

// DelCard godoc
// @Summary Delete a card by UUID
// @Description Delete a specific card identified by its UUID
// @Tags cards
// @Param id path string true "Card UUID"
// @Success 202 {string} string "Delete accepted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/cards/{id} [delete]
func (c *Controller) DelCard(w http.ResponseWriter, r *http.Request) {
	cardUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
	}

	if err = c.uc.DelCard(r.Context(), cardUUID, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte("Delete accepted")); err != nil {
		return
	}
}
