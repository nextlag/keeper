package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// getUserFromCtx - Retrieves the current user from the request context
func (c *Controller) getUserFromCtx(ctx context.Context) (entity.User, error) {
	currentUser, ok := ctx.Value(currentUserKey).(entity.User)
	if !ok {
		return entity.User{}, errs.ErrUnexpectedError
	}

	return currentUser, nil
}

// UserInfo godoc
// @Summary Get current user information
// @Description Retrieve information about the current user
// @Tags user
// @Produce json
// @Success 200 {object} entity.User
// @Failure 500 {object} response
// @Router /user/info [get]
func (c *Controller) UserInfo(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(currentUser); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
	}
}
