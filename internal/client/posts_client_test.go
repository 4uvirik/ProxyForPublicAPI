package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger/handler/slogdiscard"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func makeResponse(status int, body any) *http.Response {
	var reader io.Reader

	switch v := body.(type) {
	case string:
		reader = strings.NewReader(v)
	case []byte:
		reader = bytes.NewReader(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Sprintf("failed marshal body: %v", err))
		}
		reader = bytes.NewReader(data)
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(reader),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func TestClient_GetPost(t *testing.T) {

	type testCase struct {
		name                string
		roundErr            error
		status              int
		body                any
		expectedResp        *Resp
		expectedErrContains string
	}

	tests := []testCase{
		{
			name:                "transport error",
			roundErr:            errors.New("client error"),
			expectedResp:        nil,
			expectedErrContains: "request failed",
		},
		{
			name:                "status no 200",
			status:              500,
			body:                `{"error":"server"}`,
			expectedResp:        nil,
			expectedErrContains: "unexpected status code: 500",
		},
		{
			name:                "bad json",
			status:              200,
			body:                `not a json`,
			expectedResp:        nil,
			expectedErrContains: "request failed",
		},
		{
			name:   "Ok",
			status: 200,
			body: Resp{
				ID:    1,
				Title: "Test Title",
				Body:  "Test Body",
			},
			expectedResp: &Resp{
				ID:    1,
				Title: "Test Title",
				Body:  "Test Body",
			},
			expectedErrContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, req.Method)
					assert.Equal(t, "/posts/1", req.URL.Path)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))

					if tt.roundErr != nil {
						return nil, tt.roundErr
					}
					return makeResponse(tt.status, tt.body), nil
				}),
			}

			restyClient := resty.NewWithClient(httpClient)
			restyClient.SetBaseURL("http://example.com")
			restyClient.SetRetryCount(0)
			restyClient.SetHeader("Accept", "application/json")

			logger := slogdiscard.NewDiscardLogger()

			c := &Client{
				resty:  restyClient,
				logger: logger,
			}

			got, err := c.GetPost(context.Background(), 1)

			if tt.expectedErrContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrContains)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, got)
			}
		})
	}
}
