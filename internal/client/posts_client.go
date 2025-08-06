package client

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type Client struct {
	Client *http.Client
	URL    string
}

type Resp struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c *Client) GetPosts() ([]Resp, error) {

	url := c.URL + "/posts"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("bad NewRequest", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		slog.Error("хуй знает2", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.Error("Invalid response code", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("cant read:", err)
	}

	var response []Resp
	err = json.Unmarshal(body, &response)
	if err != nil {
		slog.Error("bad Unmarshal", err)
	}

	return response, nil
}
