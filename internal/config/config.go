package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger     LoggerConfig
	Postgres   PostgresConfig
	HTTPServer HTTPServerConfig
}

type LoggerConfig struct {
	Level    string `env:"LOG_LEVEL,required"`
	FilePath string `env:"LOG_FILE_PATH,required"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DBName   string `env:"POSTGRES_DB,required"`
	SSLMode  string `env:"POSTGRES_SSLMODE,required"`
}

type HTTPServerConfig struct {
	Addr string `env:"HTTP_SERVER_ADDR,required"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load .env file: %w", err)
	}
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse .env: %w", err)
	}
	return &config, nil
}
