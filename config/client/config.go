package client

import (
	"log"
	"log/slog"
	"sync"

	"github.com/nextlag/keeper/pkg/cleanenv"
)

type (
	// Config holds the configuration settings for the application.
	// It is populated from a YAML file and environment variables.
	Config struct {
		App          *App          `yaml:"app"`
		Server       *Server       `yaml:"server"`
		Log          *Log          `yaml:"logger"`
		SQLite       *SQLite       `yaml:"sqlite"`
		FilesStorage *FilesStorage `yaml:"files_storage"`
	}

	// App contains application-specific settings.
	App struct {
		Name string `yaml:"name" env:"APP_NAME"`
	}

	// Server contains server-related settings.
	Server struct {
		ServerURL string `yaml:"server_url" env:"SERVER_URL"`
	}

	// Log contains settings for logging.
	Log struct {
		Level slog.Level `yaml:"level" env:"LOG_LEVEL"`
	}

	// SQLite contains settings for SQLite database.
	SQLite struct {
		DSN string `yaml:"sqlite_dsn" env:"SQLITE_DSN"`
	}

	// FilesStorage contains file storage-related settings.
	FilesStorage struct {
		ServerLocation string `yaml:"server_location"`
		ClientLocation string `yaml:"client_location"`
	}
)

var (
	currentConfig *Config   // The application's configuration
	once          sync.Once // Ensures that the configuration is loaded only once
	ConfigYaml    = "./config/client/config.yaml"
)

// Load initializes and returns the configuration.
func Load() *Config {
	var err error

	once.Do(func() {
		cfg := Config{}
		err = cleanenv.ReadConfig(ConfigYaml, &cfg)
		if err != nil {
			log.Fatalf("LoadConfig: %v", err)
		}

		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatalf("LoadConfig: %v", err)
		}
		currentConfig = &cfg
	})

	return currentConfig
}
