package usecase

import (
	"github.com/google/uuid"

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

		AddCard(userPassword string, card *entity.Card)
		ShowCard(userPassword, cardID string)
		DelCard(userPassword, cardID string)
	}

	ClientRepo interface {
		MigrateDB()

		AddUser(user *entity.User) error
		UpdateUserToken(user *entity.User, token *entity.JWT) error
		DropUserToken() error
		RemoveUsers()
		UserExistsByEmail(email string) bool
		GetUserPasswordHash() string
		GetSavedAccessToken() (string, error)

		AddCard(*entity.Card) error
		SaveCards([]entity.Card) error
		LoadCards() []viewsets.CardForList
		GetCardByID(cardID uuid.UUID) (entity.Card, error)
		DelCard(cardID uuid.UUID) error
	}

	ClientAPI interface {
		Login(user *entity.User) (entity.JWT, error)
		Register(user *entity.User) error

		GetCards(accessToken string) ([]entity.Card, error)
		AddCard(accessToken string, card *entity.Card) error
		DelCard(accessToken, cardID string) error
	}
)
