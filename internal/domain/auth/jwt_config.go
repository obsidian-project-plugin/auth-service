package auth

import (
	"os"
)

func loadSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("Переменная окружения JWT_SECRET не установлена")
	}
	return secret
}
