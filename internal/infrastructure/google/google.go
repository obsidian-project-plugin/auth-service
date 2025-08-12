package google

import (
	"encoding/json"
	"fmt"
	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GoogleOAuth struct {
	config config.Config
	client *http.Client
}

func NewGoogleOAuth(cfg config.Config) *GoogleOAuth {
	return &GoogleOAuth{
		config: cfg,
		client: &http.Client{},
	}
}

func (g *GoogleOAuth) GetAuthURL(state string) string {
	v := url.Values{}
	v.Set("client_id", g.config.GoogleClientID)
	v.Set("redirect_uri", g.config.GoogleRedirectURI)
	v.Set("response_type", "code")
	v.Set("scope", strings.Join(g.config.GoogleScopes, " ")) // "email profile"
	v.Set("state", state)

	return "https://accounts.google.com/o/oauth2/v2/auth?" + v.Encode()
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (g *GoogleOAuth) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	v := url.Values{}
	v.Set("code", code)
	v.Set("client_id", g.config.GoogleClientID)
	v.Set("client_secret", g.config.GoogleClientSecret)
	v.Set("redirect_uri", g.config.GoogleRedirectURI)
	v.Set("grant_type", "authorization_code")

	resp, err := g.client.PostForm("https://oauth2.googleapis.com/token", v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось обменять код: %s", string(body))
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (g *GoogleOAuth) GetUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка получения информации пользователя: %s", string(body))
	}

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
func (g *GoogleOAuth) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, nil
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
func (g *GoogleOAuth) postForm(url string, data url.Values) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		resp, err = g.client.PostForm(url, data)
		if err != nil {
			break
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return resp, err
}
