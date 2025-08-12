package http

import (
	"github.com/gin-gonic/gin"
	"github.com/obsidian-project-plugin/auth-service/internal/config/db"
	"github.com/obsidian-project-plugin/auth-service/internal/service"
	"net/http"
	"time"
)

func RegisterRoutes(engine *gin.Engine, db *db.DbConnection) *gin.Engine {
	api := engine.Group("/api/v1")

	{
		api.GET("/health", healthCheck)
	}

	return engine
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

type Handler struct {
	authService *service.GoogleAuthService
}

func NewHandler(authService *service.GoogleAuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) InitiateGoogleAuth(c *gin.Context) {
	state, err := h.authService.GenerateState()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx := c.Request.Context()
	err = h.authService.RedisClient.Set(ctx, "google_oauth:state:"+state, "true", time.Minute*5)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	authURL := h.authService.GetAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Missing code or state"})
		return
	}

	authResult, err := h.authService.Authenticate(c.Request.Context(), code, state)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  authResult.AccessToken,
		"refresh_token": authResult.RefreshToken,
	})
}

func SetupRoutes(router *gin.Engine, handler *Handler) {
	router.GET("api/auth/google/initiate", handler.InitiateGoogleAuth)
	router.GET("api/auth/google/callback", handler.GoogleCallback)
}
