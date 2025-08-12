package google

import (
	"encoding/json"
	"fmt"
	"github.com/obsidian-project-plugin/auth-service/internal/app/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/obsidian-project-plugin/auth-service/internal/config"
)

type GoogleOAuth struct {
	cfg    *config.Config
	client *http.Client
}

func NewGoogleOAuth(cfg *config.Config) *GoogleOAuth {
	return &GoogleOAuth{cfg: cfg, client: &http.Client{}}
}

func (g *GoogleOAuth) GetAuthURL(state string) string {
	v := url.Values{}
	v.Set("client_id", g.cfg.GoogleClientID)
	v.Set("redirect_uri", g.cfg.GoogleRedirectURI)
	v.Set("response_type", "code")
	scopes := strings.Split(strings.Join(g.cfg.GoogleScopes, ","), " ")
	v.Set("scope", strings.Join(scopes, " ")) // "email profile"
	v.Set("state", state)

	return g.cfg.GoogleAuthURLPrefix + v.Encode()
}

func (g *GoogleOAuth) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	v := url.Values{}
	v.Set("code", code)
	v.Set("client_id", g.cfg.GoogleClientID)
	v.Set("client_secret", g.cfg.GoogleClientSecret)
	v.Set("redirect_uri", g.cfg.GoogleRedirectURI)
	v.Set("grant_type", g.cfg.GoogleGrantType)

	resp, err := utils.PostForm(g.cfg.GoogleTokenURL, v) // Use the function from the utils package
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе токена: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
	}

	fmt.Printf("Body: %v\n", string(body))

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе JSON: %w", err)
	}

	return &tokenResp, nil
}
