package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/nextlag/keeper/configuration"
	"github.com/nextlag/keeper/pkg/logger/l"
)

const sleepDuration = 10 * time.Second

func main() {
	cfg, err := configuration.Load()
	if err != nil {
		panic(err)
	}
	log := l.NewLogger(cfg)

	fmt.Println(cfg)
	fmt.Println(cfg.Postgres)
	fmt.Println(cfg.Logging)
	// Проверяем, запущен ли Docker
	if err = checkDocker(log); err != nil {
		log.Error("Ошибка при проверке Docker: %v", l.ErrAttr(err))
		return
	}

	// Управляем состоянием контейнера PostgreSQL
	switch {
	case !isContainerExist(cfg.Postgres.DockerContainer):
		log.Info("Контейнер PostgreSQL не существует, создаю и запускаю...")
		if err = runContainer(cfg); err != nil {
			log.Error("Ошибка при запуске контейнера PostgreSQL", l.ErrAttr(err))
			return
		}
		log.Info("Контейнер PostgreSQL создан и запущен успешно.")

	case !isContainerRunning(cfg.Postgres.DockerContainer):
		log.Info("Контейнер PostgreSQL существует, но не запущен. Запускаю контейнер...")
		if err = startContainer(cfg.Postgres.DockerContainer); err != nil {
			log.Error("Ошибка при запуске контейнера PostgreSQL", l.ErrAttr(err))
			return
		}
		log.Info("Контейнер PostgreSQL запущен успешно.")

	default:
		log.Info("Контейнер PostgreSQL уже запущен.")
	}
}

// checkDocker проверяет, запущен ли Docker, и запускает его, если это не так
func checkDocker(log *l.Logger) error {
	if err := exec.Command("docker", "info").Run(); err != nil {
		log.Info("Docker не запущен, пытаюсь запустить...")
		if err = exec.Command("open", "-a", "docker").Run(); err != nil {
			return err
		}
		log.Info("Docker запущен успешно.")

		// Ждем, пока Docker полностью запустится
		time.Sleep(sleepDuration)
	}
	return nil
}

// isContainerExist проверяет, существует ли контейнер с указанным именем
func isContainerExist(name string, log *l.Logger) bool {
	out, err := exec.Command("docker", "ps", "-a", "-q", "-f", "name="+name).Output()
	if err != nil {
		log.Error("Ошибка при проверке контейнеров", l.ErrAttr(err))
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}

// isContainerRunning проверяет, запущен ли контейнер с указанным именем
func isContainerRunning(name string, log *l.Logger) bool {
	out, err := exec.Command("docker", "ps", "-q", "-f", "name="+name).Output()
	if err != nil {
		log.Error("Ошибка при проверке запущенных контейнеров", l.ErrAttr(err))
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}

// runContainer запускает новый контейнер PostgreSQL
func runContainer(cfg *configuration.Config) error {
	return exec.Command(
		"docker",
		"run",
		"--name",
		cfg.Postgres.DockerContainer,
		"-p",
		cfg.Postgres.PostgresPort,
		"-e",
		cfg.Postgres.PostgresUser,
		"-e",
		cfg.Postgres.PostgresPass,
		"-d",
		cfg.Postgres.DockerImage,
	).Run()
}

// startContainer запускает существующий контейнер PostgreSQL
func startContainer(name string) error {
	return exec.Command("docker", "start", name).Run()
}
