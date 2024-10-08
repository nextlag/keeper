package repo

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

func (r *Repo) RemoveUsers() {
	r.db.Exec("DELETE FROM users")
}
func (r *Repo) RemoveTempUser() {
	r.db.Exec("DELETE FROM temp_users")
}

func (r *Repo) AddUser(user *entity.User) error {
	r.RemoveTempUser()
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("repo - AddUser - HashPassword - %w", err)
	}
	newUser := models.User{Email: user.Email, Password: hashedPassword}
	return r.db.Create(&newUser).Error
}

func (r *Repo) AddTempPass(user *entity.User) error {
	tempUser := models.TempUser{Email: user.Email, Password: user.Password}
	return r.db.Create(&tempUser).Error
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

func (r *Repo) DropUserToken(email string) error {
	var existedUser models.User

	result := r.db.Where("email = ?", email).First(&existedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with email %s not found", email)
		}
		return result.Error
	}

	existedUser.AccessToken = ""
	existedUser.RefreshToken = ""

	return r.db.Save(&existedUser).Error
}

func (r *Repo) GetUserPasswordHash() (string, error) {
	var existedUser models.User

	tempUser, err := r.GetTempUser()
	if err != nil {
		return "", err
	}

	result := r.db.Where("email = ?", tempUser.Email).First(&existedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user with email %s not found", tempUser.Email)
		}
		return "", fmt.Errorf("failed to query user with email %s: %v", tempUser.Email, result.Error)
	}

	return existedUser.Password, nil
}

func (r *Repo) GetSavedAccessToken() (accessToken string, err error) {
	var user models.User

	tempUser, err := r.GetTempUser()
	if err != nil {
		return "", err
	}
	result := r.db.Where("email = ?", tempUser.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user with email %s not found", tempUser.Email)
		}
		return "", result.Error
	}
	return user.AccessToken, err
}

func (r *Repo) getUserID() uint {
	var user models.User
	tempUser, err := r.GetTempUser()
	if err != nil {
		return 0
	}
	result := r.db.Where("email = ?", tempUser.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0
		}
		return 0
	}

	return user.ID
}

func (r *Repo) GetTempUser() (*models.TempUser, error) {
	var tempUser models.TempUser
	if err := r.db.First(&tempUser).Error; err != nil {
		return nil, err
	}
	return &tempUser, nil
}
