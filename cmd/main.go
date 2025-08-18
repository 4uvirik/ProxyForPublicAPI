package main

import (
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger"
	"github.com/4uvirik/ProxyForPublicAPI/internal/server"
	"log"
	"log/slog"
	"os"
)

func main() {

	yamlPath := os.Getenv("YAML_PATH")
	if yamlPath == "" {
		yamlPath = "config/config.yaml"
	}

	cfg, err := config.LoadConfig(yamlPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validate failed: %v", err)
	}

	logger, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	logger.Debug("configuration reading success", slog.Any("cfg", cfg))

	cli := client.NewClient(cfg, logger)

	h := &handlers.Handler{
		Client: cli,
		Logger: logger,
		Config: cfg,
	}

	server.Run(cfg, h, logger)

}
