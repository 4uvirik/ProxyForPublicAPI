package main

import (
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger/handler/slogpretty"
	"github.com/4uvirik/ProxyForPublicAPI/internal/server"
	"log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustNew()

	logger, err := configLogger(cfg.Logger.Level)
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

func configLogger(logLevel string) (*slog.Logger, error) {
	var log *slog.Logger

	switch logLevel {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log, nil
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
