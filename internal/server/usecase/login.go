package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
)

func (uc *UseCase) AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	return uc.repo.AddLogin(ctx, login, userID)
}
