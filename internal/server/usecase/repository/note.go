package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/usecase/repository/models"
	"github.com/nextlag/keeper/internal/utils/errs"
)

func (r *Repo) GetNotes(ctx context.Context, user entity.User) ([]entity.SecretNote, error) {
	var notesFromDB []models.Note

	err := r.db.WithContext(ctx).
		Model(&models.Note{}).
		Preload("Meta").
		Find(&notesFromDB, "user_id = ?", user.ID).Error
	if err != nil {
		return nil, err
	}

	if len(notesFromDB) == 0 {
		return nil, nil
	}

	notes := make([]entity.SecretNote, len(notesFromDB))

	for index := range notesFromDB {
		notes[index].ID = notesFromDB[index].ID
		notes[index].Name = notesFromDB[index].Name
		notes[index].Note = notesFromDB[index].Note
	}
	return notes, nil
}

func (r *Repo) AddNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		noteToDB := models.Note{
			ID:     uuid.New(),
			UserID: userID,
			Name:   note.Name,
			Note:   note.Note,
		}

		if err := r.db.WithContext(ctx).Create(&noteToDB).Error; err != nil {
			return err
		}

		note.ID = noteToDB.ID
		for index, meta := range note.Meta {
			metaForNote := models.MetaNote{
				Name:   meta.Name,
				Value:  meta.Value,
				NoteID: noteToDB.ID,
				ID:     meta.ID,
			}
			if err := tx.WithContext(ctx).Create(&metaForNote).Error; err != nil {
				return err
			}
			note.Meta[index].ID = metaForNote.ID
		}

		return nil
	})
}

func (r *Repo) IsNoteOwner(ctx context.Context, noteID, userID uuid.UUID) bool {
	var noteFromDB models.Note

	r.db.WithContext(ctx).Where("id = ?", noteID).First(&noteFromDB)

	return noteFromDB.UserID == userID
}

func (r *Repo) DelNote(ctx context.Context, noteID, userID uuid.UUID) error {
	if !r.IsNoteOwner(ctx, noteID, userID) {
		return errs.ErrWrongOwnerOrNotFound
	}

	return r.db.WithContext(ctx).Delete(&models.Note{}, noteID).Error
}

func (r *Repo) UpdateNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error {
	if !r.IsNoteOwner(ctx, note.ID, userID) {
		return errs.ErrWrongOwnerOrNotFound
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		noteToDB := models.Note{
			ID:     note.ID,
			UserID: userID,
			Name:   note.Name,
			Note:   note.Note,
		}

		if err := r.db.WithContext(ctx).Save(&noteToDB).Error; err != nil {
			return nil
		}

		return nil
	})
}
