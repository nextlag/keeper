package usecase

import (
	"log"

	"github.com/fatih/color"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
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
	if err = uc.repo.AddTempPass(user); err != nil {
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
	user, err := uc.repo.GetTempUser()
	if err != nil {
		return
	}

	if err = uc.repo.DropUserToken(user.Email); err != nil {
		color.Red("Internal error: %v", err)
		return
	}

	uc.repo.RemoveTempUser()

	color.Green("Users tokens were dropped")
}

func (uc *ClientUseCase) Sync(userPassword string) {
	if !uc.verifyPassword(userPassword) {
		return
	}
	accessToken, err := uc.repo.GetSavedAccessToken()
	if err != nil {
		color.Red("Internal error: %v", err)
		return
	}
	uc.loadLogins(accessToken)
	uc.loadCards(accessToken)
	uc.loadNotes(accessToken)
}

func (uc *ClientUseCase) verifyPassword(userPassword string) bool {
	hashPassword, err := uc.repo.GetUserPasswordHash()
	if err != nil {
		return false
	}
	if err = utils.VerifyPassword(hashPassword, userPassword); err != nil {
		color.Red("Password check status: failed")
		return false
	}
	return true
}

func (uc *ClientUseCase) GetTempPass() (string, error) {
	user, err := uc.repo.GetTempUser()
	if err != nil {
		return "", err
	}
	return user.Password, nil
}
