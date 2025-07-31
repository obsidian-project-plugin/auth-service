package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	domain "github.com/obsidian-project-plugin/auth-service/internal/domain/auth"
)

type TokenService interface {
	GenerateTokens(user domain.User, app domain.App) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(refreshToken, tokenID string) (newAccessToken string, err error)
}

type tokenService struct {
	secret    []byte
	accessTTL time.Duration
}

func NewTokenService(secret string, accessTTL time.Duration) TokenService {
	return &tokenService{
		secret:    []byte(secret),
		accessTTL: accessTTL,
	}
}

func (s *tokenService) GenerateTokens(user domain.User, app domain.App) (string, string, error) {
	now := time.Now()
	randomJTI := uuid.NewString()

	accessClaims := domain.CustomClaims{
		Sub:       user.ID,
		Email:     user.Email,
		AppID:     app.ID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			ID:        randomJTI,
		},
	}
	accessToken, err := s.sign(accessClaims)
	if err != nil {
		return "", "", err
	}

	refreshClaims := domain.CustomClaims{
		Sub:       user.ID,
		Email:     user.Email,
		AppID:     app.ID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ID:       randomJTI,
		},
	}
	refreshToken, err := s.sign(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *tokenService) RefreshAccessToken(refreshToken, tokenID string) (string, error) {
	tok, err := jwt.ParseWithClaims(refreshToken, &domain.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}
	claims, ok := tok.Claims.(*domain.CustomClaims)
	if !ok || !tok.Valid {
		return "", errors.New("invalid token")
	}
	if claims.TokenType != "refresh" {
		return "", errors.New("token is not refresh type")
	}
	if claims.ID != tokenID {
		return "", errors.New("token_id mismatch")
	}

	now := time.Now()
	newClaims := domain.CustomClaims{
		Sub:       claims.Sub,
		Email:     claims.Email,
		AppID:     claims.AppID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			ID:        tokenID,
		},
	}
	newToken, err := s.sign(newClaims)
	if err != nil {
		return "", fmt.Errorf("could not generate access token: %w", err)
	}
	return newToken, nil
}

func (s *tokenService) sign(claims domain.CustomClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}
