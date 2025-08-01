package dto

import "github.com/google/uuid"

type LoginDto struct {
	Sub   uuid.UUID `json:"sub" binding:"required"`
	Email string    `json:"email"  binding:"required"`
	AppID uuid.UUID `json:"app_id" binding:"required"`
}

type RefreshDto struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	TokenID      string `json:"token_id"      binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	CauseBy string `json:"cause_by"`
}
