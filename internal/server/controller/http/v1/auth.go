package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nextlag/keeper/internal/utils/errs"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Controller) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload *loginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.uc.SignUpUser(r.Context(), payload.Email, payload.Password)
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

	if errors.Is(err, errs.ErrWrongEmail) || errors.Is(err, errs.ErrEmailAlreadyExists) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

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
