package repo

import (
	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
	"github.com/nextlag/keeper/internal/client/usecase/viewsets"
)

func (r *Repo) LoadNotes() []viewsets.NoteForList {
	userID := r.getUserID()
	var notes []models.Note
	r.db.
		Model(&models.Note{}).
		Preload("Meta").
		Where("user_id", userID).Find(&notes)
	if len(notes) == 0 {
		return nil
	}

	notesViewSet := make([]viewsets.NoteForList, len(notes))

	for index := range notes {
		notesViewSet[index].ID = notes[index].ID
		notesViewSet[index].Name = notes[index].Name
	}

	return notesViewSet
}
