package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

type Config struct {
	Server     ServerConfig     `json:"server"`
	Logger     LoggerConfig     `json:"logger"`
	HTTPClient HTTPClientConfig `json:"HTTPClient"`
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

type HTTPClientConfig struct {
	JsonPlaceholderUrl string `env:"JSONPLACEHOLDER_URL"`
	Timeout            int    `env:"HTTP_CLIENT_TIMEOUT" envDefault:"10"`
	RetryCount         int    `env:"HTTP_RETRY_COUNT" envDefault:"3"`
	RetryWaitMs        int    `env:"HTTP_RETRY_WAIT_MS" envDefault:"200"`
	RetryMaxWaitMs     int    `env:"HTTP_RETRY_MAX_WAIT_MS" envDefault:"2000"`
}

func MustNew() *Config {

	config := Config{}

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: Error loading .env file, fallback to system ENV")
	}

	if err := env.Parse(&config); err != nil {
		log.Fatal("Failed to parse env vars into config")
	}
	return &config
}

func (cfg *Config) Validate() error {

	var errorsMsg []string

	if cfg.Server.Host == "" {
		errorsMsg = append(errorsMsg, "invalid server host")
	}

	if cfg.Server.Port == "" {
		errorsMsg = append(errorsMsg, "invalid server port")
	}

	if cfg.HTTPClient.JsonPlaceholderUrl == "" {
		errorsMsg = append(errorsMsg, "invalid client jsonPlaceholderUrl")
	}

	return fmt.Errorf(strings.Join(errorsMsg, "\n"))
}
