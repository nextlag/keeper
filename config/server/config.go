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
)

var (
	once sync.Once
	cfg  Config
)

type (
	Config struct {
		ConfigPath   string        `yaml:"config_path"`
		Network      *Network      `yaml:"network"`
		Logging      *Logging      `yaml:"logging"`
		Security     *Security     `yaml:"security"`
		PG           *PG           `yaml:"postgres"`
		Cache        *Cache        `yaml:"cache"`
		FilesStorage *FilesStorage `yaml:"files_storage"`
	}

	Network struct {
		Host  string `yaml:"host"`
		HTTPS bool   `yaml:"https"`
	}

	Logging struct {
		Level       slog.Level `yaml:"level"`
		ProjectPath string     `yaml:"project_path" env:"PROJECT_PATH"`
		LogToFile   bool       `yaml:"log_to_file" env:"LOG_TO_FILE"`
		LogPath     string     `yaml:"log_path" env:"LOG_PATH"`
	}

	Security struct {
		Domain                 string        `yaml:"domain" env:"DOMAIN"`
		AccessTokenPrivateKey  string        `yaml:"access_token_private_key" env:"ACCESS_TOKEN_PRIVATE_KEY"`
		AccessTokenPublicKey   string        `yaml:"access_token_public_key" env:"ACCESS_TOKEN_PUBLIC_KEY"`
		RefreshTokenPrivateKey string        `yaml:"refresh_token_private_key" env:"REFRESH_TOKEN_PRIVATE_KEY"`
		RefreshTokenPublicKey  string        `yaml:"refresh_token_public_key" env:"REFRESH_TOKEN_PUBLIC_KEY"`
		AccessTokenExpiresIn   time.Duration `yaml:"access_token_expired_in" env:"ACCESS_TOKEN_EXPIRED_IN"`
		RefreshTokenExpiresIn  time.Duration `yaml:"refresh_token_expired_in" env:"REFRESH_TOKEN_EXPIRED_IN"`
		AccessTokenMaxAge      int           `yaml:"access_token_maxage" env:"ACCESS_TOKEN_MAXAGE"`
		RefreshTokenMaxAge     int           `yaml:"refresh_token_maxage" env:"REFRESH_TOKEN_MAXAGE"`
	}

	PG struct {
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX"`
		DSN     string `yaml:"dsn" env:"DSN"`
	}

	Cache struct {
		DefaultExpiration int `yaml:"default_expiration" env:"DEFAULT_EXPIRATION"`
		CleanupInterval   int `yaml:"cleanup_interval" env:"CLEANUP_INTERVAL"`
	}

	FilesStorage struct {
		Location string `yaml:"location" env:"FILES_LOCATION"`
	}
)

func NewConfig() (*Config, error) {
	var (
		err     error
		homeDir string
	)
	once.Do(func() {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return
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
		flag.Var(&LogLevelValue{&cfg.Logging.Level}, "level", "Log level (debug, info, warn, error)")

		if err = env.Parse(&cfg); err != nil {
			return
		}

		flag.Parse()
	})
	return &cfg, err
}

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
