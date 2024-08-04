package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/usecase/repository/models"
	"github.com/nextlag/keeper/internal/utils/errs"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// AddLogin adds a new login entry to the database.
// It wraps the database operation in a transaction to ensure atomicity.
// Creates a new Login record with the provided entity.Login data
// and associates it with the user identified by userID.
// Also creates MetaLogin records for any metadata associated with the login entry.
// Returns an error if the operation fails.
func (r *Repo) AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) (err error) {
	return r.db.Transaction(func(tx *gorm.DB) error {
		loginToDB := models.Login{
			ID:       uuid.New(),
			UserID:   userID,
			Name:     login.Name,
			Password: login.Password,
			URI:      login.URI,
			Login:    login.Login,
		}

		if err = tx.WithContext(ctx).Create(&loginToDB).Error; err != nil {
			return l.WrapErr(err)
		}

		login.ID = loginToDB.ID
		for index, meta := range login.Meta {
			metaForLogin := models.MetaLogin{
				Name:    meta.Name,
				Value:   meta.Value,
				LoginID: loginToDB.ID,
			}
			if err = tx.WithContext(ctx).Create(&metaForLogin).Error; err != nil {
				return l.WrapErr(err)
			}
			login.Meta[index].ID = metaForLogin.ID
		}

		return nil
	})
}

// GetLogins retrieves all login entries associated with a specific user.
// Returns a slice of Login and an error if any occurs.
func (r *Repo) GetLogins(ctx context.Context, user entity.User) (logins []entity.Login, err error) {
	var loginsFromDB []models.Login

	err = r.db.WithContext(ctx).
		Model(&models.Login{}).
		Preload("Meta").
		Find(&loginsFromDB, "user_id = ?", user.ID).Error
	if err != nil {
		return nil, l.WrapErr(err)
	}

	if len(loginsFromDB) == 0 {
		return nil, nil
	}

	logins = make([]entity.Login, len(loginsFromDB))

	for index := range loginsFromDB {
		logins[index].ID = loginsFromDB[index].ID
		logins[index].Name = loginsFromDB[index].Name
		logins[index].Password = loginsFromDB[index].Password
		logins[index].URI = loginsFromDB[index].URI
		logins[index].Login = loginsFromDB[index].Login
		for metaIndex := range loginsFromDB[index].Meta {
			logins[index].Meta = append(logins[index].Meta, entity.Meta{
				ID:    loginsFromDB[index].Meta[metaIndex].ID,
				Name:  loginsFromDB[index].Meta[metaIndex].Name,
				Value: loginsFromDB[index].Meta[metaIndex].Value,
			})
		}
	}

	return
}

// IsLoginOwner checks if a specific user is the owner of a login entry.
// Returns true if the user is the owner, false otherwise.
func (r *Repo) IsLoginOwner(ctx context.Context, loginID, userID uuid.UUID) bool {
	var loginFromDB models.Login
	r.db.WithContext(ctx).Where("id = ?", loginID).First(&loginFromDB)
	return loginFromDB.UserID == userID
}

// DelLogin deletes a login entry if the user is the owner of the login.
// Returns an error if the user is not the owner or if any other issue occurs during deletion.
func (r *Repo) DelLogin(ctx context.Context, loginID, userID uuid.UUID) (err error) {
	if !r.IsLoginOwner(ctx, loginID, userID) {
		return l.WrapErr(errs.ErrWrongOwnerOrNotFound)
	}

	return l.WrapErr(r.db.WithContext(ctx).Delete(&models.Login{}, loginID).Error)
}

// UpdateLogin updates an existing login entry if the user is the owner of the login.
// Updates the login details and associated metadata.
// Returns an error if the user is not the owner or if any other issue occurs during the update.
func (r *Repo) UpdateLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	if !r.IsLoginOwner(ctx, login.ID, userID) {
		return l.WrapErr(errs.ErrWrongOwnerOrNotFound)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		loginToDB := models.Login{
			ID:       login.ID,
			Name:     login.Name,
			Password: login.Password,
			URI:      login.URI,
			Login:    login.Login,
			UserID:   userID,
		}

		if err := tx.WithContext(ctx).Save(&loginToDB).Error; err != nil {
			return l.WrapErr(err)
		}
		login.ID = loginToDB.ID
		for _, meta := range login.Meta {
			metaForLogin := models.MetaLogin{
				Name:    meta.Name,
				Value:   meta.Value,
				LoginID: loginToDB.ID,
				ID:      meta.ID,
			}
			if err := tx.WithContext(ctx).Create(&metaForLogin).Error; err != nil {
				return l.WrapErr(err)
			}
		}
		return nil
	})
}
