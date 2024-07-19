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

// GetCards handles the retrieval of cards for the current user.
// It fetches the user from the context, retrieves their cards from the use case,
// and returns them as a JSON response. If there are no cards, it returns a 204 No Content status.
func (c *Controller) GetCards(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
	}

	userCards, err := c.uc.GetCards(r.Context(), currentUser)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(userCards) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(userCards); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AddCard handles the addition of a new card for the current user.
// It fetches the user from the context, decodes the card from the request body,
// and passes it to the use case for creation. If successful, it returns the created card as a JSON response.
func (c *Controller) AddCard(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
	}

	var payloadCard *entity.Card
	if err = json.NewDecoder(r.Body).Decode(&payloadCard); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.uc.AddCard(r.Context(), payloadCard, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(payloadCard); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DelCard handles the deletion of a card by its UUID for the current user.
// It fetches the user and the card ID from the request, passes them to the use case for deletion,
// and returns a 202 Accepted status if successful.
func (c *Controller) DelCard(w http.ResponseWriter, r *http.Request) {
	cardUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
	}

	if err = c.uc.DelCard(r.Context(), cardUUID, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// UpdateCard handles the update of a card by its UUID for the current user.
// It fetches the user and the card ID from the request, decodes the updated card data from the request body,
// and passes them to the use case for updating. If successful, it returns a 202 Accepted status.
func (c *Controller) UpdateCard(w http.ResponseWriter, r *http.Request) {
	cardUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payloadCard.ID = cardUUID
	if err = c.uc.UpdateCard(r.Context(), payloadCard, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
