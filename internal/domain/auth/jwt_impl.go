package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtServiceImpl struct {
	secretKey string
}

func (s *jwtServiceImpl) GenerateJWT(userID string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiration).Unix(),
	})
	return token.SignedString([]byte(s.secretKey))
}
