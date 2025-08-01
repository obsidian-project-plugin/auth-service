package service

import (
	"context"
	"fmt"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"github.com/obsidian-project-plugin/auth-service/internal/domain"
	"golang.org/x/oauth2"
	"log"
)

type AuthService struct {
	config          config.Config
	oauth2Config    *oauth2.Config
	generateStateID func() string
}

func NewAuthService(cfg config.Config, oauth2Config *oauth2.Config, generateStateID func() string) *AuthService {
	log.Println("Создание нового auth service")
	return &AuthService{
		config:          cfg,
		oauth2Config:    oauth2Config,
		generateStateID: generateStateID,
	}
}

func (s *AuthService) GetGitHubAuthURL(redirectURL string) string {
	state := s.generateStateID()
	_ = redirectURL
	return s.oauth2Config.AuthCodeURL(state)
}

func (s *AuthService) ProcessGitHubCallback(ctx context.Context, state, code string) (string, error) {
	_ = state
	token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("Не удалось обменять токен: %w", err)
	}
	user, err := s.GetGithubUserInfo(ctx, token)
	if err != nil {
		return "", fmt.Errorf("Не удалось получить информацию о пользователе: %w", err)
	}
	log.Println(user)
	return "/", nil
}

func (s *AuthService) GetGithubUserInfo(ctx context.Context, token *oauth2.Token) (*domain.User, error) {

	return &domain.User{
		ID:       "",
		Username: "TODO",
		Email:    "TODO",
	}, nil
}
