package server

import (
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/4uvirik/ProxyForPublicAPI/internal/routes"
	"github.com/labstack/echo/v4"
)

func Run(handler *handlers.Handler) {
	e := echo.New()

	routes.InitRoutes(e, handler)

	e.Logger.Fatal(e.Start(":8080"))
}
