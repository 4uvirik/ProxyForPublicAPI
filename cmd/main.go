package main

import (
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger/handler/slogpretty"
	"github.com/4uvirik/ProxyForPublicAPI/internal/server"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
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

	httpClient := &http.Client{
		Timeout: time.Duration(cfg.HTTPClient.Timeout) * time.Second,
	}

	cli := &client.Client{
		Client:         httpClient,
		URL:            cfg.HTTPClient.JsonPlaceholderUrl,
		RetryCount:     cfg.HTTPClient.RetryCount,
		RetryWaitMs:    cfg.HTTPClient.RetryWaitMs,
		RetryMaxWaitMs: cfg.HTTPClient.RetryMaxWaitMs,
		RetryBackOff:   cfg.HTTPClient.RetryBackOff,
	}

	h := &handlers.Handler{
		Client: cli,
	}

	server.Run(cfg, h)
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
