package client

import (
	"context"
	"fmt"
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger/sl"
	"github.com/go-resty/resty/v2"
	"log/slog"
)

type Client struct {
	resty  *resty.Client
	logger *slog.Logger
}

type Resp struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

const UrlEndPointPosts = "/posts"

func NewClient(cfg *config.Config, logger *slog.Logger) *Client {
	client := resty.New().
		SetRetryCount(cfg.HTTPClient.RetryCount).
		SetRetryWaitTime(cfg.HTTPClient.RetryWait).
		SetRetryMaxWaitTime(cfg.HTTPClient.RetryMaxWait).
		SetBaseURL(cfg.HTTPClient.JsonPlaceholderUrl).
		SetHeader("Accept", "application/json").

		// Лог после завершения запроса - результат
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			logger.Info("complete HTTP request",
				slog.String("method", r.Request.Method),
				slog.String("URL", r.Request.URL),
				slog.Int("status code", r.StatusCode()),
				slog.String("requestBody", fmt.Sprintf("%s", r.Request.Body)),
				slog.String("responseBody", string(r.Body())))
			return nil
		})

	return &Client{
		resty:  client,
		logger: logger,
	}
}

func (c *Client) GetPost(ctx context.Context, postID int) (*Resp, error) {
	var response Resp
	url := fmt.Sprintf("%s/%v", UrlEndPointPosts, postID)

	resp, err := c.resty.R().
		SetContext(ctx).
		SetResult(&response).
		Get(url)
	if err != nil {
		c.logger.Error("request failed", sl.Err(err))
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		errStatus := fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		c.logger.Error("unexpected status code", sl.Err(errStatus))
		return nil, errStatus
	}

	return &response, nil
}
