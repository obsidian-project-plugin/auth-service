package auth

import "time"

type JWTService interface {
	GenerateJWT(userID string, expiration time.Duration) (string, error)
}

func NewJWTService() JWTService {
	return &jwtServiceImpl{secretKey: loadSecret()}
}
