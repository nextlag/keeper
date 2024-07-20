package usecase

import (
	"sync"

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

type UseCaseOpts func(*ClientUseCase)

func SetRepo(r ClientRepo) UseCaseOpts {
	return func(uc *ClientUseCase) {
		uc.repo = r
	}
}

func SetAPI(clientAPI ClientAPI) UseCaseOpts {
	return func(uc *ClientUseCase) {
		uc.clientAPI = clientAPI
	}
}

func SetConfig(cfg *config.Config) UseCaseOpts {
	return func(uc *ClientUseCase) {
		uc.cfg = cfg
	}
}

func (uc *ClientUseCase) InitDB() {
	uc.repo.MigrateDB()
}
