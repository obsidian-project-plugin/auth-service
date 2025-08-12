package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/obsidian-project-plugin/auth-service/internal/app/utils"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"github.com/obsidian-project-plugin/auth-service/internal/config/db"
	internalHttp "github.com/obsidian-project-plugin/auth-service/internal/presentation/http"
	"github.com/obsidian-project-plugin/auth-service/internal/telemetry/logging"
)

// Run запускает приложение: настраивает всё, запускает сервер и обрабатывает корректное завершение.
func Run() {
	cfg := loadConfig()

	ctx, stop := newSignalContext()
	defer stop()

	conn := initDB(ctx, cfg)
	defer closeDB(conn)

	router := setupRouter(conn)

	srv := createServer(cfg, router)
	startServer(srv)

	<-ctx.Done()
	shutdownServer(srv)
}

// loadConfig загружает конфигурацию приложения и инициализирует логирование.
func loadConfig() *config.Config {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		//logging.Fatal("Не удалось загрузить конфиг:", err)
		fmt.Fprintf(os.Stderr, "Не удалось загрузить конфиг: %v\n", err)
		os.Exit(1)
	}
	logging.Init(cfg) //
	return cfg
}

// newSignalContext возвращает context.Context, который отменится при получении SIGINT/SIGTERM.
func newSignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

// initDB инициализирует подключение к базе данных и возвращает объект подключения.
func initDB(ctx context.Context, cfg *config.Config) *db.DbConnection {
	conn := db.InitConnection(ctx, cfg)
	return conn
}

// closeDB закрывает подключение к базе данных, логируя возможные ошибки.
func closeDB(conn *db.DbConnection) {
	if err := conn.DB.Close(); err != nil {
		logging.Error("Ошибка при закрытии БД:", err)
	}
}

// setupRouter настраивает Gin-движок, middleware и маршруты.
func setupRouter(conn *db.DbConnection) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(utils.LoggingMiddleware(logging.Logger))

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	internalHttp.RegisterRoutes(engine, conn)
	return engine
}

// createServer создаёт HTTP-сервер на базе конфигурации и переданного handler.
func createServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.HTTPPort,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeOut) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeOut) * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// startServer запускает HTTP-сервер в отдельной горутине.
func startServer(srv *http.Server) {
	go func() {
		logging.Info("Сервер запускается на", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Fatal("Сервер упал:", err)
		}
	}()
}

// shutdownServer выполняет graceful shutdown сервера с таймаутом в 10 секунд.
func shutdownServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logging.Info("Инициализация graceful shutdown...")
	if err := srv.Shutdown(ctx); err != nil {
		logging.Error("Ошибка при завершении сервера:", err)
	} else {
		logging.Info("Сервер корректно завершён")
	}
}
