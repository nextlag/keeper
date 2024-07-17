package repository

import (
	"context"
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
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	AddLogin(ctx context.Context, login *entity.Login, userID uuid.UUID) error
	AddUser(ctx context.Context, email, hashedPassword string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email, hashedPassword string) (entity.User, error)
}

// Repo implements the Repository interface and provides methods for database operations.
type Repo struct {
	db  *gorm.DB
	log *l.Logger
}

// New creates a new Repo instance with the given database DSN and logger.
func New(dsn string, log *l.Logger) *Repo {
	attempts := 3
	for attempts > 0 {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return &Repo{
				db:  db,
				log: log,
			}
		}
		log.Info("Database: %s is not available, attempts left: %d", dsn, attempts)
		time.Sleep(time.Second)
		attempts--
	}
	log.Error("Repo - New - could not connect")
	return nil
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
