package google

import (
	"encoding/json"
	"fmt"
	"github.com/obsidian-project-plugin/auth-service/internal/app/utils"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GoogleOAuth struct {
	cfg    *config.Config
	client *http.Client
}

func NewGoogleOAuth(cfg *config.Config) *GoogleOAuth {
	return &GoogleOAuth{cfg: cfg, client: &http.Client{}}
}

func (g *GoogleOAuth) GetAuthURL(state string) string {
	values := url.Values{}
	values.Set("client_id", g.cfg.GoogleClientID)
	values.Set("redirect_uri", g.cfg.GoogleRedirectURI)
	values.Set("response_type", "code")
	scopes := strings.Split(strings.Join(g.cfg.GoogleScopes, ","), " ")
	values.Set("scope", strings.Join(scopes, " ")) // "email profile"
	values.Set("state", state)

	return g.cfg.GoogleAuthURLPrefix + values.Encode()
}

func (g *GoogleOAuth) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	values := url.Values{}
	values.Set("code", code)
	values.Set("client_id", g.cfg.GoogleClientID)
	values.Set("client_secret", g.cfg.GoogleClientSecret)
	values.Set("redirect_uri", g.cfg.GoogleRedirectURI)
	values.Set("grant_type", g.cfg.GoogleGrantType)

	resp, err := utils.PostForm(g.cfg.GoogleTokenURL, values)
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

	return &tokenResp, nil //
}
