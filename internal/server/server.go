package server

import (
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/4uvirik/ProxyForPublicAPI/internal/routes"
	"github.com/labstack/echo/v4"
	"log/slog"
)

func Run(cfg *config.Config, handler *handlers.Handler, logger *slog.Logger) {
	e := echo.New()

	routes.InitRoutes(e, handler)

	address := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Info("start server", slog.String("address", address))
	e.Logger.Fatal(e.Start(address))
}
