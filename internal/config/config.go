package config

import "github.com/spf13/viper"

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
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	Cache  CacheConfig  `mapstructure:"cache"`
	Stage  Stage        `mapstructure:"stage"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
