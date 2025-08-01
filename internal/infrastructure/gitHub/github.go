package github

import (
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"golang.org/x/oauth2"
	go_github "golang.org/x/oauth2/github"
)

func NewOAuth2Config(cfg config.GithubConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURI,
		Scopes:       cfg.Scopes,
		Endpoint:     go_github.Endpoint,
	}
}
