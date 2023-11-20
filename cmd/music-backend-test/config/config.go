package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const envfile = "./dev/.env"

type Config struct {
	LogLevel    string `long:"log-level" description:"Log level: panic, fatal, warn, debug, info" env:"LOG_LEVEL" default:"info"`
	Environment string `long:"env" description:"App environment" env:"ENV" default:"develop"`

	Debug   bool   `long:"debug" description:"Developer mode" env:"DEBUG"`
	DevMode bool   `long:"dev_mode" description:"Developer mode" env:"DEV_MODE"`
	PathLog string `long:"path_log" description:"Path log" env:"PATH_LOG" default:"stdout"`

	ApiKey string `long:"api-key" description:"Api key for authentification" env:"APIKEY" default:"apikey"`

	AppInfo struct {
		Name    string `long:"name" description:"App name" env:"APP_NAME" required:"true" default:"default app"`
		Version string `long:"version" description:"App version" env:"APP_VERSION" required:"true" default:"0.0.1"`
	}

	HttpServer struct {
		Host string `long:"http_host" description:"Host HTTP server" env:"HTTP_HOST" required:"true" default:"0.0.0.0"`
		Port int    `long:"http_port" description:"Post HTTP sever" env:"HTTP_PORT" required:"true" default:"80"`
	}

	DB struct {
		Host     string `long:"db_host" description:"Host DB" env:"DB_HOST" required:"true" default:"127.0.0.1"`
		Port     int    `long:"db_port" description:"Port DB" env:"INTERNAL_DB_PORT" required:"true" default:"5432"`
		Name     string `long:"db_name" description:"Name DB" env:"DB_NAME" required:"true" default:"db"`
		Username string `long:"db_username" description:"Username DB" env:"DB_USER" required:"true" default:"dbuser"`
		Password string `long:"db_password" description:"Password DB" env:"DB_PASS" required:"true" default:"dbpass"`
		SSLMode  string `long:"db_sslmode" description:"SSLMode DB" env:"DB_SSLMODE" required:"true" default:"disable"`
	}
}

var (
	appConfig     *Config
	appConfigOnce sync.Once
)

func newConfig() (*Config, error) {
	err := godotenv.Load(envfile)
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error parse environment variables: %v", err)
	}
	err = env.Parse(&cfg.AppInfo)
	if err != nil {
		return nil, fmt.Errorf("error parse environment variables: %v", err)
	}
	err = env.Parse(&cfg.HttpServer)
	if err != nil {
		return nil, fmt.Errorf("error parse environment variables: %v", err)
	}
	err = env.Parse(&cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("error parse environment variables: %v", err)
	}

	return &cfg, nil
}

func GetAppConfig() (*Config, error) {
	var err_ error
	appConfigOnce.Do(func() {
		config, err := newConfig()
		if err != nil {
			err_ = fmt.Errorf("can't load config: %v", err)
		}
		appConfig = config
	})

	return appConfig, err_
}
