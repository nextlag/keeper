package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
	"github.com/nextlag/keeper/pkg/logger/l"
)

var errBinaryNameNotGiven = errors.New("binary name has not given")

// AddBinary adds a new binary for the current user.
func (c *Controller) AddBinary(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	var binary entity.Binary
	if r.URL.Query().Get("name") == "" {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errBinaryNameNotGiven.Error(), http.StatusBadRequest)
		return
	}
	binary.Name = r.URL.Query().Get("name")

	_, file, err := r.FormFile("file")
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	binary.FileName = file.Filename
	if err = c.uc.AddBinary(r.Context(), &binary, file, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(binary); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetBinaries retrieves all binaries for the current user.
func (c *Controller) GetBinaries(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, errs.ErrUnexpectedError.Error(), http.StatusInternalServerError)
		return
	}

	userBinaries, err := c.uc.GetBinaries(r.Context(), currentUser)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(userBinaries) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(userBinaries); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DownloadBinary downloads a binary by its UUID.
func (c *Controller) DownloadBinary(w http.ResponseWriter, r *http.Request) {
	binaryUUID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	filePath, err := c.uc.GetUserBinary(r.Context(), &currentUser, binaryUUID)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)
	http.ServeFile(w, r, filePath)
}

// AddBinaryMeta adds metadata to a binary by its UUID.
func (c *Controller) AddBinaryMeta(w http.ResponseWriter, r *http.Request) {
	binaryUUID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	var payloadMeta []entity.Meta
	if err = json.NewDecoder(r.Body).Decode(&payloadMeta); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	binary, err := c.uc.AddBinaryMeta(r.Context(), &currentUser, binaryUUID, payloadMeta)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(binary.Meta); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DelBinary deletes a binary by its UUID.
func (c *Controller) DelBinary(w http.ResponseWriter, r *http.Request) {
	binaryUUID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	if err = c.uc.DelUserBinary(r.Context(), &currentUser, binaryUUID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
