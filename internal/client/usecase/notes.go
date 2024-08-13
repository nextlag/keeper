package usecase

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

func (uc *ClientUseCase) loadNotes(accessToken string) {
	notes, err := uc.clientAPI.GetNotes(accessToken)
	if err != nil {
		color.Red("Connection error: %v", err)
		return
	}

	if err = uc.repo.SaveNotes(notes); err != nil {
		log.Println(err)
		return
	}
	color.Green("Loaded %v notes", len(notes))
}

func (uc *ClientUseCase) AddNote(userPassword string, note *entity.SecretNote) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		log.Printf("ClientUseCase - AddNote - %v", err)
		return
	}
	uc.encryptNote(userPassword, note)

	if err = uc.clientAPI.AddNote(accessToken, note); err != nil {
		log.Printf("ClientUseCase - AddNote - %v", err)
		return
	}

	if err = uc.repo.AddNote(note); err != nil {
		log.Println(err)
		return
	}
	color.Green("Note %q added, id: %v", note.Name, note.ID)
}

func (uc *ClientUseCase) ShowNote(userPassword, noteID string) {
	if !uc.verifyPassword(userPassword) {
		return
	}
	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		color.Red(err.Error())
		return
	}

	note, err := uc.repo.GetNoteByID(noteUUID)
	if err != nil {
		color.Red(err.Error())
		return
	}

	uc.decryptNote(userPassword, &note)
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("ID: %s\nname: %s\nNote: %s\nMeta: %v\n",
		yellow(note.ID),
		yellow(note.Name),
		yellow(note.Note),
		yellow(note.Meta),
	)
}

func (uc *ClientUseCase) encryptNote(userPassword string, note *entity.SecretNote) {
	note.Note = utils.Encrypt(userPassword, note.Note)
}

func (uc *ClientUseCase) decryptNote(userPassword string, note *entity.SecretNote) {
	note.Note = utils.Decrypt(userPassword, note.Note)
}

func (uc *ClientUseCase) DelNote(userPassword, noteID string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		return
	}
	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		color.Red(err.Error())
		log.Printf("ClientUseCase - uuid.Parse - %v", err)
		return
	}

	if err = uc.repo.DelNote(noteUUID); err != nil {
		log.Printf("ClientUseCase - repo.DelNote - %v", err)
		return
	}

	if err = uc.clientAPI.DelNote(accessToken, noteID); err != nil {
		log.Printf("ClientUseCase - repo.DelNote - %v", err)
		return
	}

	color.Green("Note %q removed", noteID)
}
