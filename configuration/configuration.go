package configuration

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	cfg  Config
	once sync.Once
	err  error
)

type Config struct {
	ConfigPath string    `yaml:"config_path" env:"CONFIG_PATH" envDefault:"config.yaml"`
	Logging    *Logging  `yaml:"logging"`
	Postgres   *Postgres `yaml:"postgres"`
}

type Logging struct {
	Level       string `yaml:"level" env:"LOG_LEVEL" envDefault:"debug"`
	ProjectPath string `yaml:"project_path" env:"PROJECT_PATH"`
	LogToFile   bool   `yaml:"log_to_file" env:"LOG_TO_FILE" envDefault:"false"`
	LogPath     string `yaml:"log_path" env:"LOG_PATH" envDefault:""`
}

type Postgres struct {
	DockerImage     string `yaml:"docker_image" env:"DOCKER_IMAGE"`
	DockerContainer string `yaml:"docker_container" env:"DOCKER_CONTAINER"`
	PostgresPass    string `yaml:"postgres_pass" env:"POSTGRES_PASS"`
	PostgresUser    string `yaml:"postgres_user" env:"POSTGRES_USER"`
	PostgresPort    string `yaml:"postgres_port" env:"POSTGRES_PORT"`
}

// Load initializes the configuration by reading command line flags and environment variables.
func Load() (*Config, error) {
	once.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		envPath := fmt.Sprintf("%s/Documents/GoProjects/keeper/.env", homeDir)
		err = godotenv.Load(envPath)
		if err != nil {
			log.Fatal("error parsing .env: ", err)
		}

		// Регистрируем флаги командной строки
		// flag.StringVar(&cfg.Logging.Level, "level", "debug", "Log level (debug, info, warn, error)")

		// Получаем путь к конфигурационному файлу из переменных окружения, если указан
		if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
			cfg.ConfigPath = configPath
			err = loadConfigFromYAML()
			if err != nil {
				return
			}
		}

		// override with environment variables
		err = env.Parse(&cfg)
		if err != nil {
			return
		}

		// override with command-line flags
		flag.Parse()
	})
	return &cfg, err
}

func loadConfigFromYAML() error {
	if cfg.ConfigPath == "" {
		log.Println("the path to the configuration file is empty")
		return nil
	}

	data, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}
	return nil
}
