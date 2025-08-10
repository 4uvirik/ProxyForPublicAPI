package client

import (
	"context"
	"fmt"
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"time"
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
		SetRetryWaitTime(time.Duration(cfg.HTTPClient.RetryWaitMs)*time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(cfg.HTTPClient.RetryMaxWaitMs)*time.Millisecond).
		SetBaseURL(cfg.HTTPClient.JsonPlaceholderUrl).
		SetHeader("Accept", "application/json").
		// Лог перед запросом - отправка
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			logger.Info("send HTTP request",
				slog.String("method", r.Method),
				slog.String("URL", r.URL))
			return nil
		}).
		// Лог после ответа - результат
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			logger.Info("received HTTP response",
				slog.Int("status code", r.StatusCode()),
				slog.String("body", r.String()))
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
		c.logger.Error("request failed",
			slog.Any("error", err))
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		c.logger.Error("unexpected status code",
			slog.Int("status code", resp.StatusCode()))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return &response, nil
}
