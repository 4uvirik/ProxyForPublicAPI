package handlers

import (
	"context"
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Client *client.Client
	Logger *slog.Logger
	Config *config.Config
}

type ErrorResponse struct {
	Massage string `json:"massage"`
}

func (h *Handler) Proxy(c echo.Context) error {

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(h.Config.HTTPClient.Timeout)*time.Second)
	defer cancel()

	postIDStr := c.Param("id")
	if postIDStr == "" {
		h.Logger.Warn("no found ID in post")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Massage: "Post ID required"})
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.Logger.Warn("invalid ID format", "postID", postIDStr)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Massage: "Post ID must be number"})
	}

	post, err := h.Client.GetPost(ctx, postID)
	if err != nil {
		h.Logger.Error("failed to get post", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Massage: "failed to get post"})
	}

	return c.JSON(http.StatusOK, post)
}
