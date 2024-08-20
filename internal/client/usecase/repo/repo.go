package repo

import (
	"github.com/fatih/color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
)

type Repo struct {
	db *gorm.DB
}

func New(dbFileName string) *Repo {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	db.Logger = db.Logger.LogMode(logger.Silent)
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
		&models.TempUser{},
		&models.Card{},
		&models.MetaCard{},
		&models.Login{},
		&models.MetaLogin{},
		&models.Note{},
		&models.MetaNote{},
		&models.Binary{},
		&models.MetaBinary{},
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
