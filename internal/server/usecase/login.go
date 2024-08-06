package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
)

// AddLogin adds a new login entry for a specific user.
func (uc *UseCase) AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	return uc.repo.AddLogin(ctx, login, userID)
}

// GetLogins retrieves all login entries for a given user.
func (uc *UseCase) GetLogins(ctx context.Context, user entity.User) ([]entity.Login, error) {
	return uc.repo.GetLogins(ctx, user)
}

// DelLogin deletes a login entry for a specific user based on login ID.
func (uc *UseCase) DelLogin(ctx context.Context, loginID, userID uuid.UUID) error {
	return uc.repo.DelLogin(ctx, loginID, userID)
}

// UpdateLogin updates an existing login entry for a specific user.
func (uc *UseCase) UpdateLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error {
	return uc.repo.UpdateLogin(ctx, login, userID)
}
