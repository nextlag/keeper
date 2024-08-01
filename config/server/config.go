package server

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"

	"github.com/nextlag/keeper/config"
)

type (
	// Config holds all application settings.
	Config struct {
		ConfigPath   string        `yaml:"config_path"`
		Network      *Network      `yaml:"network"`
		Log          *Log          `yaml:"logger"`
		Security     *Security     `yaml:"security"`
		PG           *PG           `yaml:"postgres"`
		Cache        *Cache        `yaml:"cache"`
		FilesStorage *FilesStorage `yaml:"files_storage"`
	}

	// Network contains network-related settings.
	Network struct {
		Host  string `yaml:"host"`
		HTTPS bool   `yaml:"https"`
	}

	// Log contains settings for logging.
	Log struct {
		Level       slog.Level `yaml:"level" env:"LOG_LEVEL"`
		ProjectPath string     `yaml:"project_path" env:"PROJECT_PATH"`
		LogToFile   bool       `yaml:"log_to_file" env:"LOG_TO_FILE"`
		LogPath     string     `yaml:"log_path" env:"LOG_PATH"`
	}

	// Security contains security-related settings.
	Security struct {
		Domain                 string        `yaml:"domain" env:"DOMAIN"`
		AccessTokenPrivateKey  string        `yaml:"access_token_private_key" env:"ACCESS_TOKEN_PRIVATE_KEY"`
		AccessTokenPublicKey   string        `yaml:"access_token_public_key" env:"ACCESS_TOKEN_PUBLIC_KEY"`
		RefreshTokenPrivateKey string        `yaml:"refresh_token_private_key" env:"REFRESH_TOKEN_PRIVATE_KEY"`
		RefreshTokenPublicKey  string        `yaml:"refresh_token_public_key" env:"ACCESS_TOKEN_PUBLIC_KEY"`
		AccessTokenExpiresIn   time.Duration `yaml:"access_token_expired_in" env:"ACCESS_TOKEN_EXPIRED_IN"`
		RefreshTokenExpiresIn  time.Duration `yaml:"refresh_token_expired_in" env:"REFRESH_TOKEN_EXPIRED_IN"`
		AccessTokenMaxAge      int           `yaml:"access_token_maxage" env:"ACCESS_TOKEN_MAXAGE"`
		RefreshTokenMaxAge     int           `yaml:"refresh_token_maxage" env:"ACCESS_TOKEN_MAXAGE"`
	}

	// PG contains PostgreSQL-related settings.
	PG struct {
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX"`
		DSN     string `yaml:"dsn" env:"DSN"`
	}

	// Cache contains caching settings.
	Cache struct {
		DefaultExpiration time.Duration `yaml:"default_expiration" env:"DEFAULT_EXPIRATION"`
		CleanupInterval   time.Duration `yaml:"cleanup_interval" env:"CLEANUP_INTERVAL"`
	}

	// FilesStorage contains file storage settings.
	FilesStorage struct {
		Location string `yaml:"location" env:"FILES_LOCATION"`
	}
)

var (
	cfg  Config    // The application's configuration.
	once sync.Once // Ensures that the configuration is loaded only once.
)

// Load initializes and returns the configuration.
// It loads configuration from a YAML file and environment variables.
// The function is thread-safe and ensures that configuration is loaded only once.
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error getting user home directory: %v", err)
		}
		envPath := fmt.Sprintf("%s/Documents/GoProjects/keeper/.env.example", homeDir)
		if err = godotenv.Load(envPath); err != nil {
			log.Fatal("error parsing .env.example: ", err)
		}

		if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
			cfg.ConfigPath = configPath
			err = configFromYAML()
			if err != nil {
				return
			}
		}

		flag.StringVar(&cfg.Network.Host, "a", cfg.Network.Host, "Host HTTP-server")
		flag.Var(&config.LogLevelValue{Value: &cfg.Log.Level}, "level", "Log level (debug, info, warn, error)")

		if err = env.Parse(&cfg); err != nil {
			log.Fatalf("error parsing environment variables: %v", err)
		}

		flag.Parse()
	})
	return &cfg, err
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
