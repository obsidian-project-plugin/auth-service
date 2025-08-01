package main

import (
	"context"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	github "github.com/obsidian-project-plugin/auth-service/internal/infrastructure/gitHub"
	"github.com/obsidian-project-plugin/auth-service/internal/infrastructure/httpService"
	"github.com/obsidian-project-plugin/auth-service/internal/presentation"
	"github.com/obsidian-project-plugin/auth-service/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

func main() {

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	oauth2Config := github.NewOAuth2Config(cfg.Github)

	authService := service.NewAuthService(*cfg, oauth2Config, uuid.NewString)

	server := httpService.NewServer(cfg.Server)

	presentation.RegisterHandlers(server.Mux(), authService)

	go func() {
		log.Printf("Запуск сервера на %s", server.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервиса: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	log.Println("Завершение работы сервера")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Не удалось завершить работу сервера: %v", err)
	}

	log.Println("Сервер остановился")
}
