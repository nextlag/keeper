package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
)

var errServer = errors.New("got server error")

func (api *ClientAPI) Login(user *entity.User) (token entity.JWT, err error) {
	client := resty.New()
	body := fmt.Sprintf(`{"email":%q, "password":%q}`, user.Email, user.Password)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&token).
		Post(fmt.Sprintf("%s/api/v1/auth/login", api.serverURL))
	if err != nil {
		return
	}

	if resp.StatusCode() == http.StatusBadRequest || resp.StatusCode() == http.StatusInternalServerError {
		errMessage := errs.ParseServerError(resp.Body())
		color.Red("Server error: %s", errMessage)
		return token, errServer
	}
	return token, nil
}

func (api *ClientAPI) Register(user *entity.User) error {
	client := resty.New()
	body := fmt.Sprintf(`{"email":%q, "password":%q}`, user.Email, user.Password)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(user).
		Post(fmt.Sprintf("%s/api/v1/auth/register", api.serverURL))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode() == http.StatusBadRequest || resp.StatusCode() == http.StatusInternalServerError {
		errMessage := errs.ParseServerError(resp.Body())
		color.Red("Server error: %s", errMessage)
		return errServer
	}

	return nil
}
