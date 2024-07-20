package usecase

import (
	"errors"
	"sync"

	"github.com/fatih/color"

	config "github.com/nextlag/keeper/config/client"
)

type ClientUseCase struct {
	repo      ClientRepo
	clientAPI ClientAPI
	cfg       *config.Config
}

var (
	clientUseCase *ClientUseCase // pattern singleton
	once          sync.Once      // pattern singleton
)

func GetClientUseCase() *ClientUseCase {
	once.Do(func() {
		clientUseCase = &ClientUseCase{}
	})

	return clientUseCase
}

type OptsUseCase func(*ClientUseCase)

func SetRepo(r ClientRepo) OptsUseCase {
	return func(uc *ClientUseCase) {
		uc.repo = r
	}
}

func SetAPI(clientAPI ClientAPI) OptsUseCase {
	return func(uc *ClientUseCase) {
		uc.clientAPI = clientAPI
	}
}

func SetConfig(cfg *config.Config) OptsUseCase {
	return func(uc *ClientUseCase) {
		uc.cfg = cfg
	}
}

func (uc *ClientUseCase) InitDB() {
	uc.repo.MigrateDB()
}

var (
	errPasswordCheck = errors.New("wrong password")
	errToken         = errors.New("user token erroe")
)

func (uc *ClientUseCase) authorisationCheck(userPassword string) (string, error) {
	if !uc.verifyPassword(userPassword) {
		return "", errPasswordCheck
	}
	accessToken, err := uc.repo.GetSavedAccessToken()
	if err != nil || accessToken == "" {
		color.Red("User should be logged")

		return "", errToken
	}

	return accessToken, nil
}
