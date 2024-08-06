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

// AddBinary godoc
// @Summary Add a new binary
// @Description Upload a new binary for the current user
// @Tags binaries
// @Accept multipart/form-data
// @Produce json
// @Param name query string true "Binary name"
// @Param file formData file true "Binary file"
// @Success 201 {object} entity.Binary
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/binary [post]
func (c *Controller) AddBinary(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	var binary entity.Binary
	if r.URL.Query().Get("name") == "" {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errBinaryNameNotGiven), http.StatusBadRequest)
		return
	}
	binary.Name = r.URL.Query().Get("name")

	_, file, err := r.FormFile("file")
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	binary.FileName = file.Filename
	if err = c.uc.AddBinary(r.Context(), &binary, file, currentUser.ID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(binary); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// GetBinaries godoc
// @Summary Get all binaries for the current user
// @Description Retrieve all binaries uploaded by the current user
// @Tags binaries
// @Produce json
// @Success 200 {array} entity.Binary
// @Success 204 "No content"
// @Failure 500 {object} response
// @Router /user/binary [get]
func (c *Controller) GetBinaries(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	userBinaries, err := c.uc.GetBinaries(r.Context(), currentUser)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
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
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// DownloadBinary godoc
// @Summary Download a binary by UUID
// @Description Download a specific binary identified by its UUID
// @Tags binaries
// @Param id path string true "Binary UUID"
// @Success 200 {file} binary
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/binary/{id} [get]
func (c *Controller) DownloadBinary(w http.ResponseWriter, r *http.Request) {
	binaryUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	filePath, err := c.uc.GetUserBinary(r.Context(), &currentUser, binaryUUID)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)
	http.ServeFile(w, r, filePath)
}

// AddBinaryMeta godoc
// @Summary Add metadata to a binary
// @Description Add metadata to a specific binary identified by its UUID
// @Tags binaries
// @Accept json
// @Produce json
// @Param id path string true "Binary UUID"
// @Param metadata body []entity.Meta true "Metadata for the binary"
// @Success 201 {array} entity.Meta
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/binary/{id}/meta [post]
func (c *Controller) AddBinaryMeta(w http.ResponseWriter, r *http.Request) {
	binaryUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	var payloadMeta []entity.Meta
	if err = json.NewDecoder(r.Body).Decode(&payloadMeta); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	binary, err := c.uc.AddBinaryMeta(r.Context(), &currentUser, binaryUUID, payloadMeta)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(binary.Meta); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// DelBinary godoc
// @Summary Delete a binary by UUID
// @Description Delete a specific binary identified by its UUID
// @Tags binaries
// @Param id path string true "Binary UUID"
// @Success 202 {string} string "delete accepted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/binary/{id} [delete]
func (c *Controller) DelBinary(w http.ResponseWriter, r *http.Request) {
	binaryUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(errs.ErrUnexpectedError), http.StatusInternalServerError)
		return
	}

	if err = c.uc.DelUserBinary(r.Context(), &currentUser, binaryUUID); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write([]byte(jsonResponse("delete accepted"))); err != nil {
		return
	}
}
