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

	cfg := config.MustNew()

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
