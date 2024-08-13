package usecase

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

// Login authenticates a user and performs necessary actions after successful authentication.
func (uc *ClientUseCase) Login(user *entity.User) {
	token, err := uc.clientAPI.Login(user)
	if err != nil {
		color.Red("Login failed for user %s: %v", user.Email, err)
		return
	}

	if !uc.repo.UserExistsByEmail(user.Email) {
		err = uc.repo.AddUser(user)
		if err != nil {
			color.Red("Failed to add user %s to repository: %v", user.Email, err)
			return
		}
	}
	if err = uc.repo.UpdateUserToken(user, &token); err != nil {
		color.Red("Failed to update token for user %s: %v", user.Email, err)
		return
	}
	if err = uc.repo.AddTempPass(user); err != nil {
		color.Red("Failed to add temporary password for user %s: %v", user.Email, err)
		return
	}
	color.Green("Got authorization token for %q", user.Email)
}

// Register registers a new user and adds them to the repository.
func (uc *ClientUseCase) Register(user *entity.User) {
	if err := uc.clientAPI.Register(user); err != nil {
		color.Red("Registration failed for user %s: %v", user.Email, err)
		return
	}

	if err := uc.repo.AddUser(user); err != nil {
		color.Red("Failed to add registered user %s to repository: %v", user.Email, err)
		return
	}

	color.Green("User registered successfully")
	color.Green("ID: %v, email: %s", user.ID, user.Email)
}

// Logout handles the process of logging out a user by removing their temporary data.
func (uc *ClientUseCase) Logout() {
	user, err := uc.repo.GetTempUser()
	if err != nil {
		color.Red("Failed to get temporary user for logout: %v", err)
		return
	}

	if err = uc.repo.DropUserToken(user.Email); err != nil {
		color.Red("Failed to drop token for user %s: %v", user.Email, err)
		return
	}

	uc.repo.RemoveTempUser()

	color.Green("User tokens were successfully dropped")
}

// Sync synchronizes user data with the server using a valid access token.
func (uc *ClientUseCase) Sync(userPassword string) {
	if !uc.verifyPassword(userPassword) {
		color.Red("Password verification failed")
		return
	}
	accessToken, err := uc.repo.GetSavedAccessToken()
	if err != nil {
		color.Red("Failed to get saved access token: %v", err)
		return
	}
	uc.loadLogins(accessToken)
	uc.loadCards(accessToken)
	uc.loadNotes(accessToken)
	uc.loadBinaries(accessToken)
}

// verifyPassword checks if the provided password matches the stored password hash.
func (uc *ClientUseCase) verifyPassword(userPassword string) bool {
	hashPassword, err := uc.repo.GetUserPasswordHash()
	if err != nil {
		color.Red("Failed to retrieve user password hash: %v", err)
		return false
	}
	if err = utils.VerifyPassword(hashPassword, userPassword); err != nil {
		color.Red("Password check failed: %v", err)
		return false
	}
	return true
}

// GetTempPass retrieves the temporary password for the current user.
func (uc *ClientUseCase) GetTempPass() (string, error) {
	user, err := uc.repo.GetTempUser()
	if err != nil {
		return "", fmt.Errorf("failed to get temporary user: %w", err)
	}
	return user.Password, nil
}
