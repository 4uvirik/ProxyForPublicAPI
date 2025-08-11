package main

import (
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger"
	"github.com/4uvirik/ProxyForPublicAPI/internal/server"
	"log"
	"log/slog"
)

func main() {

	cfg, err := config.MustNew("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validate failed: %v", err)
	}

	logger, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("error creating logger: %s", err.Error())
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
