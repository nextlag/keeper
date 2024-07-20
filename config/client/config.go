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

	"github.com/nextlag/keeper/config"
)

type (
	Config struct {
		ConfigPath   string        `yaml:"config_path"`
		App          *App          `yaml:"app"`
		Network      *Network      `yaml:"network"`
		Log          *Log          `yaml:"logger"`
		SQLite       *SQLite       `yaml:"sqlite"`
		FilesStorage *FilesStorage `yaml:"files_storage"`
	}

	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	Network struct {
		Host string `yaml:"host" env:"HOST"`
	}

	Log struct {
		Level slog.Level `yaml:"level" env:"LOG_LEVEL"`
	}

	SQLite struct {
		DSN string `yaml:"sqlite_dsn" env:"SQLITE_DSN"`
	}

	FilesStorage struct {
		Location string `yaml:"location" env:"FILES_LOCATION"`
	}
)

var (
	cfg  Config
	once sync.Once
)

func LoadConfig() *Config {
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

		flag.StringVar(&cfg.Network.Host, "a", cfg.Network.Host, "Host HTTP-server")
		flag.Var(&config.LogLevelValue{Value: &cfg.Log.Level}, "level", "Log level (debug, info, warn, error)")

		if err = env.Parse(&cfg); err != nil {
			log.Fatalf("error parsing environment variables: %v", err)
		}

		flag.Parse()
	})
	return &cfg
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
