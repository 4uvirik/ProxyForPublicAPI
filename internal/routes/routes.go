package routes

import (
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/posts", h.Proxy)
}
