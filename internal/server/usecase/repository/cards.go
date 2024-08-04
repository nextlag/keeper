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

// GetCards retrieves all cards associated with the given user from the database.
// It loads card details and their associated meta information.
// Returns a slice of entity.Card and an error if any occurred during the operation.
func (r *Repo) GetCards(ctx context.Context, user entity.User) (cards []entity.Card, err error) {
	var cardsFromDB []models.Card

	err = r.db.WithContext(ctx).
		Model(&models.Card{}).
		Preload("Meta").
		Find(&cardsFromDB, "user_id = ?", user.ID).Error
	if err != nil {
		return nil, l.WrapErr(err)
	}

	if len(cardsFromDB) == 0 {
		return nil, nil
	}

	cards = make([]entity.Card, len(cardsFromDB))

	for index := range cardsFromDB {
		cards[index].ID = cardsFromDB[index].ID
		cards[index].Brand = cardsFromDB[index].Brand
		cards[index].CardHolderName = cardsFromDB[index].CardHolderName
		cards[index].ExpirationMonth = cardsFromDB[index].ExpirationMonth
		cards[index].ExpirationYear = cardsFromDB[index].ExpirationYear
		cards[index].Name = cardsFromDB[index].Name
		cards[index].Number = cardsFromDB[index].Number
		cards[index].SecurityCode = cardsFromDB[index].SecurityCode
	}
	return
}

// AddCard adds a new card to the database for the specified user.
// It creates a new card entry and associated meta information within a database transaction.
// If the card is successfully added, it updates the provided card entity with the new card ID.
// Returns an error if any occurred during the operation.
func (r *Repo) AddCard(ctx context.Context, card *entity.Card, userID uuid.UUID) (err error) {
	return r.db.Transaction(func(tx *gorm.DB) error {
		cardToDB := models.Card{
			ID:              uuid.New(),
			UserID:          userID,
			Name:            card.Name,
			Brand:           card.Brand,
			CardHolderName:  card.CardHolderName,
			Number:          card.Number,
			ExpirationMonth: card.ExpirationMonth,
			ExpirationYear:  card.ExpirationYear,
			SecurityCode:    card.SecurityCode,
		}

		if err = tx.WithContext(ctx).Create(&cardToDB).Error; err != nil {
			return l.WrapErr(err)
		}
		card.ID = cardToDB.ID
		for index, meta := range card.Meta {
			metaForCard := models.MetaCard{
				Name:   meta.Name,
				Value:  meta.Value,
				CardID: cardToDB.ID,
			}
			if err = tx.WithContext(ctx).Create(&metaForCard).Error; err != nil {
				return l.WrapErr(err)
			}
			card.Meta[index].ID = metaForCard.ID
		}
		return nil
	})
}

// IsCardOwner checks if the given user is the owner of the specified card.
// Returns true if the user is the owner of the card, false otherwise.
func (r *Repo) IsCardOwner(ctx context.Context, cardUUID, userID uuid.UUID) bool {
	var cardFromDB models.Card
	r.db.WithContext(ctx).Where("id = ?", cardUUID).First(&cardFromDB)
	return cardFromDB.UserID == userID
}

// DelCard deletes the specified card from the database if the user is the owner.
// It ensures the user has the right to delete the card by calling IsCardOwner.
// Returns an error if the user is not the owner or if any error occurred during deletion.
func (r *Repo) DelCard(ctx context.Context, cardUUID, userID uuid.UUID) (err error) {
	if !r.IsCardOwner(ctx, cardUUID, userID) {
		err = errs.ErrWrongOwnerOrNotFound
		return l.WrapErr(err)
	}
	return r.db.WithContext(ctx).Delete(&models.Card{}, cardUUID).Error
}

// UpdateCard updates the details of the specified card in the database if the user is the owner.
// It performs the update within a database transaction, including updating associated meta information.
// Returns an error if the user is not the owner or if any error occurred during the update.
func (r *Repo) UpdateCard(ctx context.Context, card *entity.Card, userID uuid.UUID) (err error) {
	if !r.IsCardOwner(ctx, card.ID, userID) {
		err = errs.ErrWrongOwnerOrNotFound
		return l.WrapErr(err)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		cardToDB := models.Card{
			ID:              card.ID,
			UserID:          userID,
			Name:            card.Name,
			Brand:           card.Brand,
			CardHolderName:  card.CardHolderName,
			Number:          card.Number,
			ExpirationMonth: card.ExpirationMonth,
			ExpirationYear:  card.ExpirationYear,
			SecurityCode:    card.SecurityCode,
		}
		if err = tx.WithContext(ctx).Save(&cardToDB).Error; err != nil {
			return l.WrapErr(err)
		}
		for _, meta := range card.Meta {
			metaForCard := models.MetaCard{
				Name:   meta.Name,
				Value:  meta.Value,
				CardID: cardToDB.ID,
				ID:     meta.ID,
			}
			if err = tx.WithContext(ctx).Create(&metaForCard).Error; err != nil {
				return l.WrapErr(err)
			}
		}

		return nil
	})
}
