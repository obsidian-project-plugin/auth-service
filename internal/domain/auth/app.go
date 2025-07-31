package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type App struct {
	ID               uuid.UUID `db:"id"`
	Name             string    `db:"name"`
	ClientIdentifier string    `db:"client_identifier"`
}

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
}

type CustomClaims struct {
	Sub       uuid.UUID `json:"sub"`
	Email     string    `json:"email"`
	AppID     uuid.UUID `json:"app_id"`
	TokenType string    `json:"typ"`
	jwt.RegisteredClaims
}
