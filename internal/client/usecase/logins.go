package usecase

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

func (uc *ClientUseCase) loadLogins(accessToken string) {
	logins, err := uc.clientAPI.GetLogins(accessToken)
	if err != nil {
		color.Red("Connection error: %v", err)
		return
	}

	if err = uc.repo.SaveLogins(logins); err != nil {
		log.Println(err)
		return
	}
	color.Green("Loaded %v logins", len(logins))
}

func (uc *ClientUseCase) AddLogin(userPassword string, login *entity.Login) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		return
	}
	uc.encryptLogin(userPassword, login)

	if err = uc.clientAPI.AddLogin(accessToken, login); err != nil {
		return
	}

	if err = uc.repo.AddLogin(login); err != nil {
		log.Fatal(err)
	}

	color.Green("Login %q added, id: %v", login.Name, login.ID)
}

func (uc *ClientUseCase) ShowLogin(userPassword, loginID string) {
	if !uc.verifyPassword(userPassword) {
		return
	}
	loginUUID, err := uuid.Parse(loginID)
	if err != nil {
		color.Red(err.Error())
		return
	}
	login, err := uc.repo.GetLoginByID(loginUUID)
	if err != nil {
		color.Red(err.Error())
		return
	}

	uc.decryptLogin(userPassword, &login)
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("ID: %s\nname: %s\nURI: %s\nLogin: %s\nPassword: %s\nMeta: %v\n",
		yellow(login.ID),
		yellow(login.Name),
		yellow(login.URI),
		yellow(login.Login),
		yellow(login.Password),
		yellow(login.Meta),
	)
}

func (uc *ClientUseCase) encryptLogin(userPassword string, login *entity.Login) {
	login.Login = utils.Encrypt(userPassword, login.Login)
	login.Password = utils.Encrypt(userPassword, login.Password)
}

func (uc *ClientUseCase) decryptLogin(userPassword string, login *entity.Login) {
	login.Login = utils.Decrypt(userPassword, login.Login)
	login.Password = utils.Decrypt(userPassword, login.Password)
}

func (uc *ClientUseCase) DelLogin(userPassword, loginID string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		return
	}
	loginUUID, err := uuid.Parse(loginID)
	if err != nil {
		color.Red(err.Error())
		log.Fatalf("ClientUseCase - uuid.Parse - %v", err)
	}

	if err := uc.repo.DelLogin(loginUUID); err != nil {
		log.Fatalf("ClientUseCase - repo.DelLogin - %v", err)
	}

	if err := uc.clientAPI.DelLogin(accessToken, loginID); err != nil {
		log.Fatalf("ClientUseCase - repo.DelLogin - %v", err)
	}

	color.Green("Login %q removed", loginID)
}
