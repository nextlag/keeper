package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/usecase/repository/models"
	"github.com/nextlag/keeper/internal/utils/errs"
)

// AddLogin adds a new login entry to the database.
// It wraps the database operation in a transaction to ensure atomicity.
// The method creates a new Login record with the provided entity.Login data
// and associates it with the user identified by userID.
// If successful, it also creates corresponding MetaLogin records for any metadata
// associated with the login entry.
func (r *Repo) AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		loginToDB := models.Login{
			ID:       uuid.New(),
			UserID:   userID,
			Name:     login.Name,
			Password: login.Password,
			URI:      login.URI,
			Login:    login.Login,
		}

		if err := tx.WithContext(ctx).Create(&loginToDB).Error; err != nil {
			return err
		}

		login.ID = loginToDB.ID
		for index, meta := range login.Meta {
			metaForLogin := models.MetaLogin{
				Name:    meta.Name,
				Value:   meta.Value,
				LoginID: loginToDB.ID,
			}
			if err := tx.WithContext(ctx).Create(&metaForLogin).Error; err != nil {
				return err
			}
			login.Meta[index].ID = metaForLogin.ID
		}

		return nil
	})
}

func (r *Repo) GetLogins(ctx context.Context, user entity.User) ([]entity.Login, error) {
	var loginsFromDB []models.Login

	err := r.db.WithContext(ctx).
		Model(&models.Login{}).
		Preload("Meta").
		Find(&loginsFromDB, "user_id = ?", user.ID).Error
	if err != nil {
		return nil, err
	}

	if len(loginsFromDB) == 0 {
		return nil, nil
	}

	logins := make([]entity.Login, len(loginsFromDB))

	for i, loginFromDB := range loginsFromDB {
		login := entity.Login{
			ID:       loginFromDB.ID,
			Name:     loginFromDB.Name,
			Password: loginFromDB.Password,
			URI:      loginFromDB.URI,
			Login:    loginFromDB.Login,
		}

		for _, meta := range loginFromDB.Meta {
			login.Meta = append(login.Meta, entity.Meta{
				ID:    meta.ID,
				Name:  meta.Name,
				Value: meta.Value,
			})
		}
		logins[i] = login
	}
	return logins, nil
}

func (r *Repo) IsLoginOwner(ctx context.Context, loginID, userID uuid.UUID) bool {
	var loginFromDB models.Login

	r.db.WithContext(ctx).Where("id = ?", loginID).First(&loginFromDB)

	return loginFromDB.UserID == userID
}

func (r *Repo) DelLogin(ctx context.Context, loginID, userID uuid.UUID) error {
	if !r.IsLoginOwner(ctx, loginID, userID) {
		return errs.ErrWrongOwnerOrNotFound
	}

	return r.db.WithContext(ctx).Delete(&models.Login{}, loginID).Error
}

func (r *Repo) UpdateLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	if !r.IsLoginOwner(ctx, login.ID, userID) {
		return errs.ErrWrongOwnerOrNotFound
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
			return err
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
				return err
			}
		}
		return nil
	})
}
