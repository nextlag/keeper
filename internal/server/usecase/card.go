package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
)

// GetCards retrieves all cards for a given user.
func (uc *UseCase) GetCards(ctx context.Context, user entity.User) ([]entity.Card, error) {
	return uc.repo.GetCards(ctx, user)
}

// AddCard adds a new card for a specific user.
func (uc *UseCase) AddCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error {
	return uc.repo.AddCard(ctx, card, userID)
}

// DelCard deletes a card for a specific user based on card UUID.
func (uc *UseCase) DelCard(ctx context.Context, cardUUID, userID uuid.UUID) error {
	return uc.repo.DelCard(ctx, cardUUID, userID)
}

// UpdateCard updates an existing card for a specific user.
func (uc *UseCase) UpdateCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error {
	return uc.repo.UpdateCard(ctx, card, userID)
}
