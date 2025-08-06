package client

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

type Client struct {
	Client         *http.Client
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
			time.Sleep(wait)
		}
	}

	return nil, fmt.Errorf("fail after %d retries", c.RetryCount)
}
