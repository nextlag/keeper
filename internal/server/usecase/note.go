package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
)

// GetNotes retrieves all secret notes for a specific user.
func (uc *UseCase) GetNotes(ctx context.Context, user entity.User) ([]entity.SecretNote, error) {
	return uc.repo.GetNotes(ctx, user)
}

// AddNote adds a new secret note for a specific user.
func (uc *UseCase) AddNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error {
	return uc.repo.AddNote(ctx, note, userID)
}

// DelNote deletes a secret note for a specific user based on note ID.
func (uc *UseCase) DelNote(ctx context.Context, noteID, userID uuid.UUID) error {
	return uc.repo.DelNote(ctx, noteID, userID)
}

// UpdateNote updates an existing secret note for a specific user.
func (uc *UseCase) UpdateNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error {
	return uc.repo.UpdateNote(ctx, note, userID)
}
