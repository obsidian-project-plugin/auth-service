package router

import (
	"github.com/gin-gonic/gin"
	"github.com/obsidian-project-plugin/auth-service/internal/config/db"
	authHandlers "github.com/obsidian-project-plugin/auth-service/internal/presentation/http/auth"
	"github.com/obsidian-project-plugin/auth-service/internal/service"
)

type AuthRouter struct {
	DbConnection *db.DbConnection
	TokenSvc     service.TokenService
}

func NewAuthRouter(dbConn *db.DbConnection, tokenSvc service.TokenService) *AuthRouter {
	return &AuthRouter{DbConnection: dbConn, TokenSvc: tokenSvc}
}

func (authRouter *AuthRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	authGroup := routerGroup.Group("auth")
	{

		{
			authHandler := authHandlers.Init(authRouter.TokenSvc)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
		}
	}
}
