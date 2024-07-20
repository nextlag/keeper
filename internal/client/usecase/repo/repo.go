package repo

import (
	"github.com/fatih/color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
)

type Repo struct {
	db *gorm.DB
}

func New(dbFileName string) *Repo {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		color.Red("Load error %s", err.Error())
	}

	return &Repo{
		db: db,
	}
}

func (r *Repo) MigrateDB() {
	tables := []interface{}{
		&models.User{},
		&models.Card{},
		&models.Login{},
		&models.Note{},
	}
	var err error
	for _, table := range tables {

		if err = r.db.Migrator().DropTable(table); err != nil {
			color.Red("Init error %s", err.Error())
		}
		if err = r.db.Migrator().CreateTable(table); err != nil {
			color.Red("Init error %s", err.Error())
		}
	}

	color.Green("Initialization status: success")
	color.Green("You can use keeper")
}
