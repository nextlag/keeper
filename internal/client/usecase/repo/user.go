package repo

import (
	"fmt"

	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

func (r *Repo) RemoveUsers() {
	r.db.Exec("DELETE FROM users")
}

func (r *Repo) AddUser(user *entity.User) error {
	r.RemoveUsers()
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("Repo - AddUser - HashPassword - %w", err)
	}

	newUser := models.User{
		Email:    user.Email,
		Password: hashedPassword,
	}

	return r.db.Create(&newUser).Error
}

func (r *Repo) UpdateUserToken(user *entity.User, token *entity.JWT) error {
	var existedUser models.User

	r.db.Where("email", user.Email).First(&existedUser)
	existedUser.AccessToken = token.AccessToken
	existedUser.RefreshToken = token.RefreshToken

	return r.db.Save(&existedUser).Error
}

func (r *Repo) UserExistsByEmail(email string) bool {
	var user models.User

	r.db.Where("email = ?", email).First(&user)

	return user.ID != 0
}

func (r *Repo) DropUserToken() error {
	var existedUser models.User

	r.db.First(&existedUser)
	existedUser.AccessToken = ""
	existedUser.RefreshToken = ""

	return r.db.Save(&existedUser).Error
}

func (r *Repo) GetUserPasswordHash() string {
	var existedUser models.User
	r.db.First(&existedUser)
	return existedUser.Password
}

func (r *Repo) GetSavedAccessToken() (accessToken string, err error) {
	var user models.User
	err = r.db.First(&user).Error

	return user.AccessToken, err
}

func (r *Repo) getUserID() uint {
	var user models.User
	r.db.First(&user)

	return user.ID
}
