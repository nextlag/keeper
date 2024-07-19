package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
)

func (uc *UseCase) AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	return uc.repo.AddLogin(ctx, login, userID)
}

func (uc *UseCase) GetLogins(ctx context.Context, user entity.User) ([]entity.Login, error) {
	return uc.repo.GetLogins(ctx, user)
}

func (uc *UseCase) DelLogin(ctx context.Context, loginID, userID uuid.UUID) error {
	return uc.repo.DelLogin(ctx, loginID, userID)
}

func (uc *UseCase) UpdateLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	return uc.repo.UpdateLogin(ctx, login, userID)
}
