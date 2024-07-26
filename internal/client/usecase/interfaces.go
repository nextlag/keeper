package usecase

import (
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
	"github.com/nextlag/keeper/internal/client/usecase/viewsets"
	"github.com/nextlag/keeper/internal/entity"
)

type (
	// Client - use cases.
	Client interface {
		InitDB()

		Register(user *entity.User)
		Login(user *entity.User)
		Logout()
		GetTempPass() (string, error)

		AddCard(userPassword string, card *entity.Card)
		ShowCard(userPassword, cardID string)
		DelCard(userPassword, cardID string)

		AddLogin(userPassword string, login *entity.Login)
		ShowLogin(userPassword, loginID string)
		DelLogin(userPassword, loginID string)

		AddNote(userPassword string, note *entity.SecretNote)
		ShowNote(userPassword, noteID string)
		DelNote(userPassword, noteID string)
	}

	ClientRepo interface {
		MigrateDB()

		AddUser(user *entity.User) error
		AddTempPass(user *entity.User) error
		UpdateUserToken(user *entity.User, token *entity.JWT) error
		DropUserToken(email string) error
		RemoveUsers()
		RemoveTempUser()
		UserExistsByEmail(email string) bool
		GetUserPasswordHash() (string, error)
		GetSavedAccessToken() (string, error)
		GetTempUser() (*models.TempUser, error)

		AddLogin(*entity.Login) error
		SaveLogins([]entity.Login) error
		LoadLogins() []viewsets.LoginForList
		GetLoginByID(loginID uuid.UUID) (entity.Login, error)
		DelLogin(loginID uuid.UUID) error

		AddCard(*entity.Card) error
		SaveCards([]entity.Card) error
		LoadCards() []viewsets.CardForList
		GetCardByID(cardID uuid.UUID) (entity.Card, error)
		DelCard(cardID uuid.UUID) error

		SaveNotes([]entity.SecretNote) error
		AddNote(*entity.SecretNote) error
		LoadNotes() []viewsets.NoteForList
		GetNoteByID(notedID uuid.UUID) (entity.SecretNote, error)
		DelNote(noteID uuid.UUID) error

		LoadBinaries() []viewsets.BinaryForList
	}

	ClientAPI interface {
		Login(user *entity.User) (entity.JWT, error)
		Register(user *entity.User) error

		AddCard(accessToken string, card *entity.Card) error
		GetCards(accessToken string) ([]entity.Card, error)
		DelCard(accessToken, cardID string) error

		AddLogin(accessToken string, login *entity.Login) error
		GetLogins(accessToken string) ([]entity.Login, error)
		DelLogin(accessToken, loginID string) error

		GetNotes(accessToken string) ([]entity.SecretNote, error)
		AddNote(accessToken string, note *entity.SecretNote) error
		DelNote(accessToken, noteID string) error
	}
)
