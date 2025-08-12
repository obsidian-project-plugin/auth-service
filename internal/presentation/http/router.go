package http

import (
	"github.com/gin-gonic/gin"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"github.com/obsidian-project-plugin/auth-service/internal/config/db"
	"github.com/obsidian-project-plugin/auth-service/internal/presentation/http/auth/router"
	"github.com/obsidian-project-plugin/auth-service/internal/service"
	"net/http"
	"time"
)

func RegisterRoutes(engine *gin.Engine, dbConn *db.DbConnection) *gin.Engine {
	api := engine.Group("/api/v1")

	{
		api.GET("/health", healthCheck)
	}

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	tokenSvc := service.NewTokenService(
		cfg.Server.Token,
		15*time.Minute,
	)

	authRouter := router.NewAuthRouter(dbConn, tokenSvc)
	authRouter.RegisterRoutes(api)

	return engine
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
