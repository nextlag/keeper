package client

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type (
	// Config holds the configuration settings for the application.
	// It is populated from a YAML file and environment variables.
	Config struct {
		ConfigPath   string        `yaml:"config_path"`
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
	cfg  Config    // The application's configuration
	once sync.Once // Ensures that the configuration is loaded only once
)

// Load initializes and returns the configuration.
// It loads configuration from a YAML file and environment variables.
// The function is thread-safe and ensures that configuration is loaded only once.
func Load() *Config {
	once.Do(func() {
		var err error
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error getting user home directory: %v", err)
		}

		envPath := fmt.Sprintf("%s/Documents/GoProjects/keeper/.env.example", homeDir)
		if err = godotenv.Load(envPath); err != nil {
			log.Fatalf("error loading .env.example: %v", err)
		}

		if configPath := os.Getenv("CLIENT_CONFIG_PATH"); configPath != "" {
			cfg.ConfigPath = configPath
			if err = configFromYAML(); err != nil {
				log.Fatalf("error loading config from YAML: %v", err)
			}
		}

		flag.StringVar(&cfg.Server.ServerURL, "a", cfg.Server.ServerURL, "Host HTTP-server")

		if err = env.Parse(&cfg); err != nil {
			log.Fatalf("error parsing environment variables: %v", err)
		}

		flag.Parse()
	})
	return &cfg
}

// configFromYAML loads configuration from a YAML file.
// It reads the YAML file specified by ConfigPath and unmarshals it into the Config structure.
// Returns an error if the file cannot be read or parsed.
func configFromYAML() error {
	if cfg.ConfigPath == "" {
		log.Println("the path to the config file is empty")
		return nil
	}

	data, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	return nil
}
