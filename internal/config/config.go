package config

import (
	"errors"
	"fmt"
	"github.com/obsidian-project-plugin/auth-service/internal/telemetry/logging"
	"github.com/spf13/viper"
	"os"
)

type ServerConfig struct {
	HTTPPort     string `mapstructure:"http_port"`
	Token        string `mapstructure:"token"`
	ReadTimeOut  int    `mapstructure:"read_time_out"`
	WriteTimeOut int    `mapstructure:"write_time_out"`
}

type DBConfig struct {
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	Login              string `mapstructure:"login"`
	Password           string `mapstructure:"password"`
	DbName             string `mapstructure:"db_name"`
	MaxPoolConnections int    `mapstructure:"max_pool_connections"`
}

type CacheConfig struct {
	RedisAddr string `mapstructure:"redis_addr"`
}

type Stage struct {
	IsDev       bool   `mapstructure:"is_dev"`
	LogFilePath string `mapstructure:"log_file_path"`
}

type Config struct {
	Server              ServerConfig `mapstructure:"server"`
	Github              GithubConfig `mapstructure:"github"`
	DB                  DBConfig     `mapstructure:"db"`
	Cache               CacheConfig  `mapstructure:"cache"`
	Stage               Stage        `mapstructure:"stage"`
	GoogleClientID      string       `mapstructure:"google_client_id"`
	GoogleTokenURL      string       `mapstructure:"google_token_url"`
	GoogleGrantType     string       `mapstructure:"google_grant_type"`
	GoogleRedirectURI   string       `mapstructure:"google_redirect_uri"`
	GoogleScopes        []string     `mapstructure:"google_scopes"`
	GoogleClientSecret  string       `mapstructure:"google_client_secret"`
	GoogleAuthURLPrefix string       `mapstructure:"google_auth_url_prefix"`
}
type GithubConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string
	RedirectURI  string   `mapstructure:"redirect_uri"`
	Scopes       []string `mapstructure:"scopes"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {

		logging.Error(fmt.Sprintf("Ошибка чтения конфиг файла: %s", err.Error()))
		return nil, errors.New("Ошибка чтения конфиг файла: " + err.Error())
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {

		logging.Error(fmt.Sprintf("Ошибка разгрузки конфигурации: %s", err.Error()))
		return nil, errors.New("Ошибка разгрузки конфигурации: " + err.Error())
	}

	cfg.Github.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	if cfg.Github.ClientSecret == "" {

		logging.Error("GITHUB_CLIENT_SECRET переменная окружения не задана")
		return nil, errors.New("GITHUB_CLIENT_SECRET переменная окружения не задана")
	}

	return &cfg, nil

}
