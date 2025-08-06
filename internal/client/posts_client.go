package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"time"
)

type Client struct {
	Client         *http.Client
	Logger         *slog.Logger
	URL            string
	RetryCount     int
	RetryWaitMs    int
	RetryMaxWaitMs int
	RetryBackOff   string
}

type Resp struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

const UrlEndPointPosts = "/posts"

func (c *Client) GetPosts() ([]Resp, error) {

	url := c.URL + UrlEndPointPosts
	var response []Resp

	c.Logger.Info("sending request to API",
		slog.String("url", url),
		slog.Int("retry count", c.RetryCount),
		slog.String("retry backOff", c.RetryBackOff),
	)

	for i := 1; i <= c.RetryCount; i++ {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		resp, err := c.Client.Do(req)
		if err == nil {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("error read body: %w", err)
			}

			err = json.Unmarshal(body, &response)
			if err != nil {
				return nil, fmt.Errorf("error unmarshal body: %w", err)
			}

			return response, nil
		}

		c.Logger.Warn("Request failed. Retry if possible",
			slog.Int("attempt", i),
			slog.Int("retry count", c.RetryCount),
			slog.Any("error", err),
		)

		if i < c.RetryCount {
			var wait time.Duration
			if c.RetryBackOff == "exponential" {
				waitMs := math.Min(
					float64(c.RetryWaitMs)*math.Pow(2, float64(i)),
					float64(c.RetryMaxWaitMs),
				)
				wait = time.Duration(waitMs) * time.Millisecond
			} else {
				wait = time.Duration(c.RetryWaitMs) * time.Millisecond
			}

			c.Logger.Debug("Wait before next try",
				slog.Duration("wait", wait))

			time.Sleep(wait)
		}
	}

	c.Logger.Error("request failed after all tries",
		slog.Int("retry count", c.RetryCount),
	)

	return nil, fmt.Errorf("fail after %d retries", c.RetryCount)
}
