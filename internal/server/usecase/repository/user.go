package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/usecase/repository/models"
	"github.com/nextlag/keeper/internal/utils"
	"github.com/nextlag/keeper/internal/utils/errs"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// AddUser inserts a new user into the database.
func (r *Repo) AddUser(ctx context.Context, email, hashedPassword string) (user entity.User, err error) {
	newUser := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	result := r.db.WithContext(ctx).Create(&newUser)
	if result.Error == nil {
		user.ID = newUser.ID
		user.Email = newUser.Email
		return
	}

	switch errs.ParsePostgresErr(result.Error).Code {
	case "23505":
		r.log.Error("error", l.ErrAttr(result.Error))
		err = errs.ErrEmailAlreadyExists
		return
	default:
		err = fmt.Errorf("%s", l.ErrAttr(result.Error))
		return
	}
}

// GetUserByID retrieves a user from the database by ID.
func (r *Repo) GetUserByID(ctx context.Context, id string) (user entity.User, err error) {
	var userFromDB models.User
	r.db.WithContext(ctx).Where("id = ?", id).First(&userFromDB)

	if userFromDB.ID == uuid.Nil {
		err = errs.ErrWrongCredentials
		return
	}

	user.ID = userFromDB.ID
	user.Email = userFromDB.Email
	return
}

// GetUserByEmail retrieves a user from the database by their email address.
// It uses the provided email to query the database and fetch the corresponding user record.
// If no user is found with the given email, it returns ErrWrongCredentials.
// It then verifies the provided hashed password against the stored password hash in the database.
// If the passwords do not match, it also returns ErrWrongCredentials.
// Upon successful retrieval and password verification, it populates and returns the user entity with ID and email.
func (r *Repo) GetUserByEmail(ctx context.Context, email, hashedPassword string) (user entity.User, err error) {
	var userFromDB models.User
	r.db.WithContext(ctx).Where("email = ?", email).First(&userFromDB)

	if userFromDB.ID == uuid.Nil {
		err = errs.ErrWrongCredentials
		return
	}

	if err = utils.VerifyPassword(userFromDB.Password, hashedPassword); err != nil {
		err = errs.ErrWrongCredentials
		return
	}

	user.ID = userFromDB.ID
	user.Email = userFromDB.Email
	return
}
