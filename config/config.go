package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Server     ServerConfig     `json:"server"`
	Logger     LoggerConfig     `json:"logger"`
	HTTPClient HTTPClientConfig `json:"HTTPClient"`
}

type ServerConfig struct {
	Host         string `env:"SERVER_HOST"`
	Port         string `env:"SERVER_PORT"`
	Env          string `env:"SERVER_ENV"`
	ReadTimeout  int    `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `env:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  int    `env:"SERVER_IDLE_TIMEOUT"`
}

type LoggerConfig struct {
	Level  string `env:"LOG_LEVEL"`
	Format string `env:"LOG_FORMAT"`
}

type HTTPClientConfig struct {
	JsonPlaceholderUrl string `env:"JSONPLACEHOLDER_URL"`
	Timeout            int    `env:"HTTP_CLIENT_TIMEOUT"`
	RetryCount         int    `env:"HTTP_RETRY_COUNT"`
	RetryWaitMs        int    `env:"HTTP_RETRY_WAIT_MS"`
	RetryMaxWaitMs     int    `env:"HTTP_RETRY_MAX_WAIT_MS"`
	RetryBachOff       string `env:"HTTP_RETRY_BACKOFF"`
}

func MustNew() *Config {

	config := Config{}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := env.Parse(&config); err != nil {
		log.Fatal("Failed to parse env vars into config")
	}
	return &config
}
