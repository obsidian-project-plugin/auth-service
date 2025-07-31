// internal/presentation/http/auth/handler.go
package auth

import (
	"github.com/obsidian-project-plugin/auth-service/internal/domain/auth"
	"github.com/obsidian-project-plugin/auth-service/internal/presentation/http/common"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/obsidian-project-plugin/auth-service/internal/presentation/http/auth/dto"
	"github.com/obsidian-project-plugin/auth-service/internal/service"
)

type Handler struct {
	TokenService service.TokenService
}

func NewTokenHandler(tokenSvc service.TokenService) *Handler {
	return &Handler{TokenService: tokenSvc}
}

/*func (tokenHandler *Handler) Login(c *gin.Context) {
	var req dto.LoginDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.CreateError("invalid request body", err))
		return
	}

	access, refresh, err := tokenHandler.TokenService.GenerateTokens(
		auth.User{ID: req.Sub, Email: req.Email},
		auth.App{ID: req.AppID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CreateError("could not generate tokens", err))
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}*/

func (tokenHandler *Handler) Login(c *gin.Context) {
	var req dto.LoginDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.CreateError("invalid request body", err))
		return
	}

	user := auth.User{
		ID:    req.Sub,
		Email: req.Email,
	}
	app := auth.App{
		ID: req.AppID,
	}

	access, refresh, err := tokenHandler.TokenService.GenerateTokens(user, app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CreateError("could not generate tokens", err))
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (tokenHandler *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{"invalid request body"})
		return
	}

	newAccess, err := tokenHandler.TokenService.RefreshAccessToken(req.RefreshToken, req.TokenID)
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshResponse{
		AccessToken:  newAccess,
		RefreshToken: req.RefreshToken,
	})
}
