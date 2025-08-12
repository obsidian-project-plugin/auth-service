package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JWTService interface {
	GenerateJWT(userID string, expiration time.Duration) (string, error)
}

type JWTServiceImpl struct {
	secretKey string
}

func NewJWTService() JWTService {

	secret := os.Getenv("jwt_secret")
	if secret == "" {
		panic("Переменная окружения JWT_SECRET не установлена")
	}
	return &JWTServiceImpl{secretKey: secret}
}

func (s *JWTServiceImpl) GenerateJWT(userID string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secretKey))
	return tokenString, err
}
