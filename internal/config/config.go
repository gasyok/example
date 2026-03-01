package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		HttpPort    string `envconfig:"APP_CONFIG_PORT" default:"8080"`
		LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
		DatabaseURL string `envconfig:"DATABASE_URL"`
		Postgres    Postgres
	}

	Postgres struct {
		Host     string `envconfig:"APP_CONFIG_POSTGRES_HOST" default:"localhost"`
		Port     string `envconfig:"APP_CONFIG_POSTGRES_PORT" default:"5432"`
		UserName string `envconfig:"APP_CONFIG_POSTGRES_USERNAME"`
		Password string `envconfig:"APP_CONFIG_POSTGRES_PASSWORD"`
		DBName   string `envconfig:"APP_CONFIG_POSTGRES_DBNAME"`
	}
)

func GetConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	return &cfg, nil
}
