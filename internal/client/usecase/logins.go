package usecase

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

// loadLogins loads logins using the API and saves them to the repository.
func (uc *ClientUseCase) loadLogins(accessToken string) {
	logins, err := uc.clientAPI.GetLogins(accessToken)
	if err != nil {
		color.Red("Error fetching logins with access token %s: %v", accessToken, err)
		return
	}

	if err = uc.repo.SaveLogins(logins); err != nil {
		color.Red("Error saving logins to repository: %v", err)
		return
	}
	color.Green("Loaded %v logins successfully", len(logins))
}

// AddLogin adds a new login for the user.
func (uc *ClientUseCase) AddLogin(userPassword string, login *entity.Login) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization check failed for user with provided password: %v", err)
		return
	}
	uc.encryptLogin(userPassword, login)

	if err = uc.clientAPI.AddLogin(accessToken, login); err != nil {
		color.Red("Error adding login %q with access token %s: %v", login.Name, accessToken, err)
		return
	}

	if err = uc.repo.AddLogin(login); err != nil {
		color.Red("Error adding login %q to repository: %v", login.Name, err)
		return
	}
	color.Green("Login %q added successfully, ID: %v", login.Name, login.ID)
}

// ShowLogin displays the login by its ID.
func (uc *ClientUseCase) ShowLogin(userPassword, loginID string) {
	if !uc.verifyPassword(userPassword) {
		color.Red("Password verification failed")
		return
	}
	loginUUID, err := uuid.Parse(loginID)
	if err != nil {
		color.Red("Error parsing login ID %s: %v", loginID, err)
		return
	}

	login, err := uc.repo.GetLoginByID(loginUUID)
	if err != nil {
		color.Red("Error fetching login with ID %s from repository: %v", loginID, err)
		return
	}

	uc.decryptLogin(userPassword, &login)
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("ID: %s\nName: %s\nURI: %s\nLogin: %s\nPassword: %s\nMeta: %v\n",
		yellow(login.ID),
		yellow(login.Name),
		yellow(login.URI),
		yellow(login.Login),
		yellow(login.Password),
		yellow(login.Meta),
	)
}

// encryptLogin encrypts the login and password using the user's password.
func (uc *ClientUseCase) encryptLogin(userPassword string, login *entity.Login) {
	login.Login = utils.Encrypt(userPassword, login.Login)
	login.Password = utils.Encrypt(userPassword, login.Password)
}

// decryptLogin decrypts the login and password using the user's password.
func (uc *ClientUseCase) decryptLogin(userPassword string, login *entity.Login) {
	login.Login = utils.Decrypt(userPassword, login.Login)
	login.Password = utils.Decrypt(userPassword, login.Password)
}

// DelLogin deletes a login by its ID.
func (uc *ClientUseCase) DelLogin(userPassword, loginID string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization check failed for user with provided password: %v", err)
		return
	}
	loginUUID, err := uuid.Parse(loginID)
	if err != nil {
		color.Red("Error parsing login ID %s: %v", loginID, err)
		return
	}

	if err = uc.repo.DelLogin(loginUUID); err != nil {
		color.Red("Error deleting login with ID %s from repository: %v", loginID, err)
		return
	}

	if err = uc.clientAPI.DelLogin(accessToken, loginID); err != nil {
		color.Red("Error deleting login %s with access token %s: %v", loginID, accessToken, err)
		return
	}
	color.Green("Login %q removed successfully", loginID)
}
