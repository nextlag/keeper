package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
)

func (uc *UseCase) GetCards(ctx context.Context, user entity.User) ([]entity.Card, error) {
	return uc.repo.GetCards(ctx, user)
}

func (uc *UseCase) AddCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error {
	return uc.repo.AddCard(ctx, card, userID)
}

func (uc *UseCase) DelCard(ctx context.Context, cardUUID, userID uuid.UUID) error {
	return uc.repo.DelCard(ctx, cardUUID, userID)
}

func (uc *UseCase) UpdateCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error {
	return uc.repo.UpdateCard(ctx, card, userID)
}
