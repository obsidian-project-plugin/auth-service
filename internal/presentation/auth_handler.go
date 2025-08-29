package presentation

import (
	"github.com/obsidian-project-plugin/auth-service/internal/service"
	"log"
	"net/http"
)

func RegisterHandlers(mux *http.ServeMux, authService *service.AuthService) {
	mux.HandleFunc("/api/auth/github/login", func(w http.ResponseWriter, r *http.Request) {
		githubLoginHandler(w, r, authService)
	})
	mux.HandleFunc("/api/auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		githubCallbackHandler(w, r, authService)
	})
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request, authService *service.AuthService) {
	redirectURL := r.URL.Query().Get("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}
	authURL := authService.GetGitHubAuthURL(redirectURL)
	log.Printf("Перенаправить на: %s", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request, authService *service.AuthService) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	redirectURL, err := authService.ProcessGitHubCallback(r.Context(), state, code)
	if err != nil {
		log.Printf("Обработка ошибок обратного вызова GitHub: %v", err)
		http.Error(w, "Не удалось обработать обратный вызов GitHub", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешный вход в систему. Перенаправление на: %s", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
