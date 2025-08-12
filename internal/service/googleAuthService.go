package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"github.com/obsidian-project-plugin/auth-service/internal/domain/auth"
	"github.com/obsidian-project-plugin/auth-service/internal/infrastructure/google"
	"github.com/obsidian-project-plugin/auth-service/internal/infrastructure/redis"
	"math/rand"
	"os"
	"time"
)

type GoogleAuthService struct {
	config      config.Config
	googleOAuth *google.GoogleOAuth
	RedisClient *redis.RedisClient
	db          *sql.DB
	jwtSecret   string
}

func NewGoogleAuthService(cfg config.Config, googleOAuth *google.GoogleOAuth, redisClient *redis.RedisClient, db *sql.DB) *GoogleAuthService {
	return &GoogleAuthService{
		config:      cfg,
		googleOAuth: googleOAuth,
		RedisClient: redisClient,
		db:          db,
		jwtSecret:   os.Getenv("JWT_SECRET"),
	}
}

func (s *GoogleAuthService) GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (s *GoogleAuthService) GetAuthURL(state string) string {
	return s.googleOAuth.GetAuthURL(state)
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

func (s *GoogleAuthService) Authenticate(ctx context.Context, code string, state string) (*AuthResult, error) {

	storedState, err := s.RedisClient.Get(ctx, "google_oauth:state:"+state)
	if err != nil {
		return nil, fmt.Errorf("invalid state")
	}
	if storedState == "" {
		return nil, fmt.Errorf("state не найден")
	}

	err = s.RedisClient.Delete(ctx, "google_oauth:state:"+state)
	if err != nil {
		return nil, fmt.Errorf("ошибка в удаление state: %w", err)
	}

	tokenResp, err := s.googleOAuth.ExchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("Ошибка обмена токена на код: %w", err)
	}

	userInfo, err := s.googleOAuth.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в получении информации о пользователи: %w", err)
	}

	user, err := s.findByGoogleSub(ctx, userInfo.Sub)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в поиске пользователя: %w", err)
	}

	if user == nil {
		user = &auth.User{
			ID:            uuid.New().String(),
			GoogleSub:     userInfo.Sub,
			Email:         userInfo.Email,
			Name:          userInfo.Name,
			GivenName:     userInfo.GivenName,
			FamilyName:    userInfo.FamilyName,
			EmailVerified: userInfo.EmailVerified,
		}
		err = s.createUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("ошибка в создание user: %w", err)
		}
	}

	accessToken, refreshToken, err := s.generateJWTTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка в создание JWT tokens: %w", err)
	}

	err = s.RedisClient.Set(ctx, "refresh_token:"+refreshToken, user.ID, time.Hour*24*30)
	if err != nil {
		return nil, fmt.Errorf("не удалось сохранить токен обновления: %w", err)
	}

	return &AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *GoogleAuthService) generateJWTTokens(userID string) (string, string, error) {

	accessToken, err := s.generateJWT(userID, time.Minute*15)
	if err != nil {
		return "", "", fmt.Errorf("ошибка в генерации access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateJWT(userID, time.Hour*24*30)
	if err != nil {
		return "", "", fmt.Errorf("ошибка в генерации refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *GoogleAuthService) generateJWT(userID string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("не удалось подписать токен: %w", err)
	}

	return signedToken, nil
}

func (s *GoogleAuthService) findByGoogleSub(ctx context.Context, googleSub string) (*auth.User, error) {
	query := `SELECT id, google_sub, email, name, given_name, family_name, picture, email_verified, created_at, updated_at FROM users WHERE google_sub = $1`

	row := s.db.QueryRowContext(ctx, query, googleSub)
	user := &auth.User{}

	err := row.Scan(&user.ID, &user.GoogleSub, &user.Email, &user.Name, &user.GivenName, &user.FamilyName, &user.EmailVerified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *GoogleAuthService) createUser(ctx context.Context, user *auth.User) error {
	query := `
		INSERT INTO users (id, google_sub, email, name, given_name, family_name, picture, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	_, err := s.db.ExecContext(ctx, query, user.ID, user.GoogleSub, user.Email, user.Name, user.GivenName, user.FamilyName, user.EmailVerified)
	return err
}
