package handlers

import (
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type Handler struct {
	Client *client.Client
	Logger *slog.Logger
}

func (h *Handler) Proxy(c echo.Context) error {

	posts, err := h.Client.GetPosts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed GetPost"})
	}
	return c.JSON(http.StatusOK, posts)
}
