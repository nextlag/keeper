package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/usecase/repository/models"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// Repository interface defines methods for interacting with the database.
type Repository interface {
	DBHealthCheck() error
	AddUser(ctx context.Context, email, hashedPassword string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email, hashedPassword string) (entity.User, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)

	GetLogins(ctx context.Context, user entity.User) ([]entity.Login, error)
	AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error
	DelLogin(ctx context.Context, loginID, userID uuid.UUID) error
	UpdateLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error
	IsLoginOwner(ctx context.Context, loginID, userID uuid.UUID) bool

	GetCards(ctx context.Context, user entity.User) ([]entity.Card, error)
	AddCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error
	DelCard(ctx context.Context, cardUUID, userID uuid.UUID) error
	UpdateCard(ctx context.Context, card *entity.Card, userID uuid.UUID) error
	IsCardOwner(ctx context.Context, cardUUID, userID uuid.UUID) bool

	GetNotes(ctx context.Context, user entity.User) ([]entity.SecretNote, error)
	AddNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error
	DelNote(ctx context.Context, noteID, userID uuid.UUID) error
	UpdateNote(ctx context.Context, note *entity.SecretNote, userID uuid.UUID) error
	IsNoteOwner(ctx context.Context, noteID, userID uuid.UUID) bool

	GetBinaries(ctx context.Context, user entity.User) ([]entity.Binary, error)
	AddBinary(ctx context.Context, binary *entity.Binary, userID uuid.UUID) error
	GetBinary(ctx context.Context, binaryID, userID uuid.UUID) (*entity.Binary, error)
	DelUserBinary(ctx context.Context, currentUser *entity.User, binaryUUID uuid.UUID) error
	AddBinaryMeta(ctx context.Context, currentUser *entity.User, binaryUUID uuid.UUID, meta []entity.Meta) (*entity.Binary, error)
}

// Repo implements the Repository interface and provides methods for database operations.
type Repo struct {
	db  *gorm.DB
	log *l.Logger
}

// New создает новый экземпляр Repo с указанным DSN базы данных и логгером.
func New(dsn string, log *l.Logger) (*Repo, error) {
	if log == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	attempts := 3
	for attempts > 0 {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Warn("database: %s is not available, attempts left: %d", dsn, attempts)
			time.Sleep(time.Second)
			attempts--
			continue
		}
		return &Repo{
			db:  db,
			log: log,
		}, nil
	}
	return nil, fmt.Errorf("could not connect to the database after several attempts")
}

// Migrate performs database schema migration for all registered models.
func (r *Repo) Migrate() {
	tables := []interface{}{
		&models.User{},
		&models.Card{},
		&models.MetaCard{},
		&models.Login{},
		&models.MetaLogin{},
		&models.Note{},
		&models.MetaNote{},
		&models.Binary{},
		&models.MetaBinary{},
	}

	if err := r.db.AutoMigrate(tables...); err != nil {
		r.log.Error("Migrate", l.ErrAttr(err))
		panic(err)
	}

	r.log.Debug("Migrate success")
}

// DBHealthCheck verifies the database connection by pinging the database.
func (r *Repo) DBHealthCheck() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// ShutDown closes the database connection.
func (r *Repo) ShutDown() {
	db, err := r.db.DB()
	if err != nil {
		r.log.Error("error", l.ErrAttr(err))
	}

	db.Close()
	r.log.Debug("db connection closed")
}
