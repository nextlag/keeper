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

// AddNote adds a new note for the current user based on the provided JSON payload.
// If successful, it returns the created note in JSON format with a 202 Accepted status.
// It handles errors by returning appropriate HTTP status codes and messages.
func (c *Controller) AddNote(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	var payloadNote entity.SecretNote

	if err = json.NewDecoder(r.Body).Decode(&payloadNote); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.uc.AddNote(r.Context(), &payloadNote, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(payloadNote); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetNotes retrieves all notes associated with the current user.
// If successful, it returns a list of notes in JSON format. If no notes are found, it returns a 204 No Content status.
// It handles errors by returning appropriate HTTP status codes and messages.
func (c *Controller) GetNotes(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	userNotes, err := c.uc.GetNotes(r.Context(), currentUser)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateNote updates an existing note identified by its UUID with the provided JSON payload.
// It returns a 202 Accepted status if the update is successful. Errors are handled by returning appropriate HTTP status codes and messages.
func (c *Controller) UpdateNote(w http.ResponseWriter, r *http.Request) {
	noteUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "noteUUID", noteUUID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	var payloadNote entity.SecretNote

	if err = json.NewDecoder(r.Body).Decode(&payloadNote); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payloadNote.ID = noteUUID

	if err = c.uc.UpdateNote(r.Context(), &payloadNote, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte("Update accepted")); err != nil {
		return
	}
}

// DelNote deletes a specific note by its UUID if it belongs to the current user.
// It returns a 202 Accepted status if the note is successfully deleted. If there are any errors, appropriate HTTP status codes and messages are returned.
func (c *Controller) DelNote(w http.ResponseWriter, r *http.Request) {
	noteUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err), "noteUUID", noteUUID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	if err = c.uc.DelNote(r.Context(), noteUUID, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte("Delete accepted")); err != nil {
		return
	}
}
