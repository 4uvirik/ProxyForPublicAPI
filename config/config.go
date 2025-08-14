package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Logger     LoggerConfig     `yaml:"logger"`
	HTTPClient HTTPClientConfig `yaml:"http_client"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST"`
	Port string `yaml:"port" env:"SERVER_PORT"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" envDefault:"info"`
}

type HTTPClientConfig struct {
	JsonPlaceholderUrl string        `yaml:"json_placeholder_url" env:"JSON_PLACEHOLDER_URL"`
	Timeout            time.Duration `yaml:"timeout" env:"HTTP_CLIENT_TIMEOUT" envDefault:"10s"`
	RetryCount         int           `yaml:"retry_count" env:"HTTP_RETRY_COUNT" envDefault:"3"`
	RetryWait          time.Duration `yaml:"retry_wait" env:"HTTP_RETRY_WAIT" envDefault:"200ms"`
	RetryMaxWait       time.Duration `yaml:"retry_max_wait" env:"HTTP_RETRY_MAX_WAIT" envDefault:"2000ms"`
}

func loadFromEnv(config *Config) error {

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("cant load env file: %w", err)
	}
	if err := env.Parse(config); err != nil {
		return fmt.Errorf("error to parse env in config: %w", err)
	}
	return nil
}

func loadFromYaml(path string, config *Config) error {
	date, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read yaml file: %w", err)
	}

	if err := yaml.Unmarshal(date, config); err != nil {
		return fmt.Errorf("cant yaml unmarshal: %w", err)
	}
	return nil
}

func LoadConfig(yamlPath string) (*Config, error) {

	config := &Config{}

	if err := loadFromEnv(config); err == nil {
		return config, nil
	}

	if err := loadFromYaml(yamlPath, config); err != nil {
		return nil, err
	}

	return config, nil
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

	if len(errorsMsg) > 0 {
		return fmt.Errorf(strings.Join(errorsMsg, ": "))
	}
	return nil
}
