package usecase

import "github.com/nextlag/keeper/internal/entity"

type (
	// Client - use cases.
	Client interface {
		InitDB()

		Register(user *entity.User)
		Login(user *entity.User)
		Logout()
	}

	ClientRepo interface {
		MigrateDB()

		AddUser(user *entity.User) error
		UpdateUserToken(user *entity.User, token *entity.JWT) error
		DropUserToken() error
		RemoveUsers()
		UserExistsByEmail(email string) bool
	}

	ClientAPI interface {
		Login(user *entity.User) (entity.JWT, error)
		Register(user *entity.User) error
	}
)
