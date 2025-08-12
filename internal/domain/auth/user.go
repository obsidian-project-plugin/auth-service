package auth

type User struct {
	ID            string `json:"id"`             // Внутренний ID пользователя в вашей системе
	GoogleSub     string `json:"google_sub"`     // ID пользователя в Google (sub claim)
	Email         string `json:"email"`          // Email пользователя
	Name          string `json:"name"`           // Полное имя пользователя
	GivenName     string `json:"given_name"`     // Имя пользователя
	FamilyName    string `json:"family_name"`    // Фамилия пользователя
	EmailVerified bool   `json:"email_verified"` // Подтвержден ли email пользователя
}
