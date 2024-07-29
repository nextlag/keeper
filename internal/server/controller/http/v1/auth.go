package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/nextlag/keeper/internal/utils/errs"
	"github.com/nextlag/keeper/pkg/logger/l"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpUser handles the HTTP request to sign up a new user.
// It decodes the JSON payload from the request body into a loginPayload struct.
// If decoding fails, it responds with a HTTP 400 Bad Request error.
// It then calls the SignUpUser method of the use case (uc) to create a new user.
// If an error occurs during sign-up, it responds with a HTTP 500 Internal Server Error and logs the error.
// If the error is due to a wrong email format or a duplicate email, it responds with a HTTP 400 Bad Request error.
// If sign-up is successful, it responds with a HTTP 201 Created status and JSON representation of the user.
func (c *Controller) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload *loginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.uc.SignUpUser(r.Context(), payload.Email, payload.Password)
	if errors.Is(err, errs.ErrWrongEmail) || errors.Is(err, errs.ErrEmailAlreadyExists) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SignInUser handles the HTTP request to sign in a user.
// It decodes the JSON payload from the request body into a loginPayload struct.
// If decoding fails, it responds with a HTTP 400 Bad Request error.
// It then calls the SignInUser method of the use case (uc) to authenticate the user and generate JWT tokens.
// If authentication fails due to wrong credentials, it responds with a HTTP 400 Bad Request error.
// For other errors during sign-in, it responds with a HTTP 500 Internal Server Error and logs the error.
// If sign-in is successful, it sets HTTP cookies for access_token, refresh_token, and a logged_in flag.
// It responds with a HTTP 200 OK status and JSON representation of the JWT tokens.
func (c *Controller) SignInUser(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwtToken, err := c.uc.SignInUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, errs.ErrWrongCredentials) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, jsonError(err), http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken.AccessToken,
		MaxAge:   jwtToken.AccessTokenMaxAge,
		Path:     "/",
		Domain:   jwtToken.Domain,
		HttpOnly: true,
		Secure:   c.cfg.Network.HTTPS,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    jwtToken.RefreshToken,
		MaxAge:   jwtToken.RefreshTokenMaxAge,
		Path:     "/",
		Domain:   jwtToken.Domain,
		HttpOnly: true,
		Secure:   c.cfg.Network.HTTPS,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		MaxAge:   jwtToken.AccessTokenMaxAge,
		Path:     "/",
		Domain:   jwtToken.Domain,
		HttpOnly: false, // Не HTTPOnly, чтобы JavaScript мог читать
		Secure:   c.cfg.Network.HTTPS,
		SameSite: http.SameSiteNoneMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(jwtToken)
	if err != nil {
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// RefreshAccessToken - handler for refreshing the access token using the provided refresh token in cookies.
// This method reads the refresh token from the "refresh_token" cookie, attempts to refresh the access token,
// and sets the new access token and "logged_in" cookie. If the refresh token is not found or invalid,
// an error response is returned.
func (c *Controller) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		c.log.Error("RefreshAccessToken: refresh token not found", l.ErrAttr(err))
		http.Error(w, "refresh token has not been found", http.StatusBadRequest)
		return
	}

	jwt, err := c.uc.RefreshAccessToken(ctx, refreshToken.Value)
	if err != nil {
		c.log.Error("RefreshAccessToken: unable to refresh access token", l.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwt.AccessToken,
		MaxAge:   jwt.AccessTokenMaxAge,
		Path:     "/",
		Domain:   jwt.Domain,
		HttpOnly: true,
		Secure:   c.cfg.Network.HTTPS,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		MaxAge:   jwt.AccessTokenMaxAge,
		Path:     "/",
		Domain:   jwt.Domain,
		HttpOnly: false,
		Secure:   c.cfg.Network.HTTPS,
		SameSite: http.SameSiteNoneMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(jwt); err != nil {
		c.log.Error("RefreshAccessToken: unable to encode response", l.ErrAttr(err))
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}

// LogoutUser - handler for logging out the current user by clearing the access, refresh, and logged_in cookies.
// This method invalidates the user's session by setting the cookies to expire immediately.
func (c *Controller) LogoutUser(w http.ResponseWriter, _ *http.Request) {
	domainName := c.uc.GetDomainName()

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   domainName,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Secure:   c.cfg.Network.HTTPS,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   domainName,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Secure:   c.cfg.Network.HTTPS,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "logged_in",
		Value:    "",
		Path:     "/",
		Domain:   domainName,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: false,
		Secure:   c.cfg.Network.HTTPS,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(jsonResponse("success"))); err != nil {
		return
	}
}
