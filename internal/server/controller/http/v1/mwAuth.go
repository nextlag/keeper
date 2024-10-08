package v1

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/nextlag/keeper/pkg/logger/l"
)

type contextKey string

const (
	currentUserKey contextKey = "currentUser"
	errorKey       contextKey = "errorKey"
)

// MwAuth returns middleware to validate the access token.
// The token can be passed either in the Authorization header or in the "access_token" cookie.
// If the token is valid, information about the user is stored in the context of the request.
// If the token is missing or invalid, the status 401 Unauthorized is returned.
func (c *Controller) MwAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var accessToken string
			accessTokenFromCookie, err := r.Cookie("access_token")
			authorizationHeader := strings.Fields(r.Header.Get("Authorization"))

			if len(authorizationHeader) > 1 && authorizationHeader[0] == "Bearer" {
				accessToken = authorizationHeader[1]
			} else if err == nil {
				accessToken = accessTokenFromCookie.Value
			}

			if accessToken == "" {
				http.Error(w, jsonError(errors.New("you are not logged in")), http.StatusUnauthorized)
				return
			}

			user, err := c.uc.CheckAccessToken(r.Context(), accessToken)
			if err != nil {
				c.log.Error("error", l.ErrAttr(err))
				http.Error(w, jsonError(err), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), currentUserKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
