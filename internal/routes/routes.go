package routes

import (
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/labstack/echo/v4"
)

const PostIDParam = "/:id"

func InitRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET(client.UrlEndPointPosts+PostIDParam, h.Proxy)
}
