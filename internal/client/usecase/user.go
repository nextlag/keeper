package usecase

import (
	"log"

	"github.com/fatih/color"

	"github.com/nextlag/keeper/internal/entity"
)

func (uc *ClientUseCase) Login(user *entity.User) {
	token, err := uc.clientAPI.Login(user)
	if err != nil {
		return
	}

	if !uc.repo.UserExistsByEmail(user.Email) {
		err = uc.repo.AddUser(user)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err = uc.repo.UpdateUserToken(user, &token); err != nil {
		log.Fatal(err)
	}
	color.Green("Got authorization token for %q", user.Email)
}

func (uc *ClientUseCase) Register(user *entity.User) {
	if err := uc.clientAPI.Register(user); err != nil {
		return
	}

	if err := uc.repo.AddUser(user); err != nil {
		color.Red("Internal error: %v", err)

		return
	}

	color.Green("User registered")
	color.Green("ID: %v", user.ID)
	color.Green("Email: %s", user.Email)
}

func (uc *ClientUseCase) Logout() {
	if err := uc.repo.DropUserToken(); err != nil {
		color.Red("Internal error: %v", err)

		return
	}

	color.Green("Users tokens were dropped")
}
