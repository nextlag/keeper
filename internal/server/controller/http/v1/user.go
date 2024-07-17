package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
)

// getUserFromCtx - Retrieves the current user from the request context
func (c *Controller) getUserFromCtx(ctx context.Context) (entity.User, error) {
	currentUser, ok := ctx.Value("currentUser").(entity.User)
	if !ok {
		return entity.User{}, errs.ErrUnexpectedError
	}

	return currentUser, nil
}

// UserInfo - handler for obtaining information about the current user
func (c *Controller) UserInfo(w http.ResponseWriter, r *http.Request) {
	currentUser, err := c.getUserFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(currentUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
