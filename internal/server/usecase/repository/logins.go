package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/usecase/repository/models"
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
