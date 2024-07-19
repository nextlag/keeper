package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
)

func (uc *UseCase) GetNotes(ctx context.Context, user entity.User) ([]entity.SecretNote, error) {
	return uc.repo.GetNotes(ctx, user)
}

func (uc *UseCase) AddNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error {
	return uc.repo.AddNote(ctx, note, userID)
}

func (uc *UseCase) DelNote(ctx context.Context, noteID, userID uuid.UUID) error {
	return uc.repo.DelNote(ctx, noteID, userID)
}

func (uc *UseCase) UpdateNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error {
	return uc.repo.UpdateNote(ctx, note, userID)
}
