package api

import (
	"github.com/nextlag/keeper/internal/entity"
)

const loginsEndpoint = "api/v1/user/logins"

func (api *ClientAPI) GetLogins(accessToken string) (logins []entity.Login, err error) {
	if err := api.getEntities(&logins, accessToken, loginsEndpoint); err != nil {
		return nil, err
	}

	return logins, nil
}

func (api *ClientAPI) AddLogin(accessToken string, login *entity.Login) error {
	return api.addEntity(login, accessToken, loginsEndpoint)
}

func (api *ClientAPI) DelLogin(accessToken, loginID string) error {
	return api.delEntity(accessToken, loginsEndpoint, loginID)
}
