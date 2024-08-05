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

// AddNote godoc
// @Summary Add a new note
// @Description Upload a new note for the current user
// @Tags notes
// @Accept json
// @Produce json
// @Param note body entity.SecretNote true "Note data"
// @Success 202 {object} entity.SecretNote
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/notes [post]
func (c *Controller) AddNote(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	var payloadNote entity.SecretNote

	if err = json.NewDecoder(r.Body).Decode(&payloadNote); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	if err = c.uc.AddNote(r.Context(), &payloadNote, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(payloadNote); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// GetNotes godoc
// @Summary Get all notes for the current user
// @Description Retrieve all notes for the current user
// @Tags notes
// @Produce json
// @Success 200 {array} entity.SecretNote
// @Success 204 "No content"
// @Failure 500 {object} response
// @Router /user/notes [get]
func (c *Controller) GetNotes(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	userNotes, err := c.uc.GetNotes(r.Context(), currentUser)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if len(userNotes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(userNotes); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// UpdateNote godoc
// @Summary Update a note by UUID
// @Description Update a specific note identified by its UUID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note UUID"
// @Param note body entity.SecretNote true "Updated note data"
// @Success 202 {string} string "Update accepted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/notes/{id} [patch]
func (c *Controller) UpdateNote(w http.ResponseWriter, r *http.Request) {
	noteUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "noteUUID", noteUUID)
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	var payloadNote entity.SecretNote

	if err = json.NewDecoder(r.Body).Decode(&payloadNote); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	payloadNote.ID = noteUUID

	if err = c.uc.UpdateNote(r.Context(), &payloadNote, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte(jsonResponse("update accepted"))); err != nil {
		return
	}
}

// DelNote godoc
// @Summary Delete a note by UUID
// @Description Delete a specific note identified by its UUID
// @Tags notes
// @Param id path string true "Note UUID"
// @Success 202 {string} string "Delete accepted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/notes/{id} [delete]
func (c *Controller) DelNote(w http.ResponseWriter, r *http.Request) {
	noteUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "noteUUID", noteUUID)
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	if err = c.uc.DelNote(r.Context(), noteUUID, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte(jsonResponse("delete accepted"))); err != nil {
		return
	}
}
