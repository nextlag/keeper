package usecase

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

// loadNotes loads notes using the API and saves them to the repository.
func (uc *ClientUseCase) loadNotes(accessToken string) {
	notes, err := uc.clientAPI.GetNotes(accessToken)
	if err != nil {
		color.Red("Error while fetching notes with access token %s: %v", accessToken, err)
		return
	}

	if err = uc.repo.SaveNotes(notes); err != nil {
		color.Red("Error while saving notes to repository: %v", err)
		return
	}
	color.Green("Loaded %v notes successfully", len(notes))
}

// AddNote adds a new note for the user.
func (uc *ClientUseCase) AddNote(userPassword string, note *entity.SecretNote) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization check failed for user with provided password: %v", err)
		return
	}
	uc.encryptNote(userPassword, note)

	if err = uc.clientAPI.AddNote(accessToken, note); err != nil {
		color.Red("Error while adding note %q with access token %s: %v", note.Name, accessToken, err)
		return
	}

	if err = uc.repo.AddNote(note); err != nil {
		color.Red("Error while adding note %q to repository: %v", note.Name, err)
		return
	}
	color.Green("Note %q added successfully, ID: %v", note.Name, note.ID)
}

// ShowNote displays a note by its ID.
func (uc *ClientUseCase) ShowNote(userPassword, noteID string) {
	if !uc.verifyPassword(userPassword) {
		color.Red("Password verification failed")
		return
	}
	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		color.Red("Error parsing note ID %s: %v", noteID, err)
		return
	}

	note, err := uc.repo.GetNoteByID(noteUUID)
	if err != nil {
		color.Red("Error while fetching note with ID %s from repository: %v", noteID, err)
		return
	}

	uc.decryptNote(userPassword, &note)
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("ID: %s\nName: %s\nNote: %s\nMeta: %v\n",
		yellow(note.ID),
		yellow(note.Name),
		yellow(note.Note),
		yellow(note.Meta),
	)
}

// encryptNote encrypts the note using the user's password.
func (uc *ClientUseCase) encryptNote(userPassword string, note *entity.SecretNote) {
	note.Note = utils.Encrypt(userPassword, note.Note)
}

// decryptNote decrypts the note using the user's password.
func (uc *ClientUseCase) decryptNote(userPassword string, note *entity.SecretNote) {
	note.Note = utils.Decrypt(userPassword, note.Note)
}

// DelNote deletes a note by its ID.
func (uc *ClientUseCase) DelNote(userPassword, noteID string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization check failed for user with provided password: %v", err)
		return
	}
	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		color.Red("Error parsing note ID %s: %v", noteID, err)
		return
	}

	if err = uc.repo.DelNote(noteUUID); err != nil {
		color.Red("Error while deleting note with ID %s from repository: %v", noteID, err)
		return
	}

	if err = uc.clientAPI.DelNote(accessToken, noteID); err != nil {
		color.Red("Error while deleting note %s with access token %s: %v", noteID, accessToken, err)
		return
	}

	color.Green("Note %q removed successfully", noteID)
}
