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

// SignUpUser godoc
// @Summary Sign up a new user
// @Description Register a new user and generate initial JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body loginPayload true "Registration credentials"
// @Success 201 {object} entity.User
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /auth/register [post]
func (c *Controller) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload *loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	user, err := c.uc.SignUpUser(r.Context(), payload.Email, payload.Password)
	if errors.Is(err, errs.ErrWrongEmail) || errors.Is(err, errs.ErrEmailAlreadyExists) {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// SignInUser godoc
// @Summary Sign in a user
// @Description Authenticate a user and generate JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body loginPayload true "Login credentials"
// @Success 200 {object} entity.JWT
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /auth/login [post]
func (c *Controller) SignInUser(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	jwtToken, err := c.uc.SignInUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, errs.ErrWrongCredentials) {
			c.log.Error("error", l.ErrAttr(err))
			http.Error(w, jsonError(err), http.StatusBadRequest)
		} else {
			c.log.Error("error", l.ErrAttr(err))
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
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// RefreshAccessToken godoc
// @Summary Refresh JWT access token
// @Description Refresh the JWT access token using the refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} entity.JWT
// @Failure 400 {object} response
// @Router /auth/refresh [post]
func (c *Controller) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, "refresh token has not been found", http.StatusBadRequest)
		return
	}

	jwt, err := c.uc.RefreshAccessToken(ctx, refreshToken.Value)
	if err != nil {
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusBadRequest)
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
		c.log.Error("error", l.ErrAttr(err))
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// LogoutUser godoc
// @Summary Log out the user
// @Description Clear JWT tokens and user session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response
// @Failure 500 {object} response
// @Router /auth/logout [post]
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
	if _, err := w.Write([]byte(jsonResponse("logout success"))); err != nil {
		return
	}
}
