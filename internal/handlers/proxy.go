package handlers

import (
	"context"
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger/sl"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type PostsClient interface {
	GetPost(ctx context.Context, postID int) (*client.Resp, error)
}
type Handler struct {
	Client PostsClient
	Logger *slog.Logger
	Config *config.Config
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Proxy(c echo.Context) error {

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.Config.HTTPClient.Timeout)
	defer cancel()

	postIDStr := c.Param("id")
	if postIDStr == "" {
		h.Logger.Warn("no found ID in post")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Post ID required"})
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.Logger.Warn("invalid ID format", "postID", postIDStr)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Post ID must be number"})
	}

	post, err := h.Client.GetPost(ctx, postID)
	if err != nil {
		h.Logger.Error("failed to get post", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to get post"})
	}

	return c.JSON(http.StatusOK, post)
}
