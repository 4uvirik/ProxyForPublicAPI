package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/4uvirik/ProxyForPublicAPI/config"
	"github.com/4uvirik/ProxyForPublicAPI/internal/client"
	"github.com/4uvirik/ProxyForPublicAPI/internal/logger/handler/slogdiscard"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockPostGetter struct {
	getPostFunc func(ctx context.Context, postID int) (*client.Resp, error)
}

func (m *mockPostGetter) GetPost(ctx context.Context, postID int) (*client.Resp, error) {
	return m.getPostFunc(ctx, postID)
}

func TestHandler_Proxy(t *testing.T) {

	type testCase struct {
		name         string
		idParam      string
		mockResponse *client.Resp
		mockError    error
		expectedCode int
		expectedBody any
	}

	tests := []testCase{
		{
			name:         "no ID",
			idParam:      "",
			expectedCode: http.StatusBadRequest,
			expectedBody: ErrorResponse{Message: "Post ID required"},
		},
		{
			name:         "invalid ID",
			idParam:      "asshole",
			expectedCode: http.StatusBadRequest,
			expectedBody: ErrorResponse{Message: "Post ID must be number"},
		},
		{
			name:         "client error",
			idParam:      "1",
			mockError:    errors.New("service unavailable"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: ErrorResponse{Message: "failed to get post"},
		},
		{
			name:    "Ok",
			idParam: "1",
			mockResponse: &client.Resp{
				ID:    1,
				Title: "Test Title",
				Body:  "Test Body",
			},
			expectedCode: http.StatusOK,
			expectedBody: client.Resp{
				ID:    1,
				Title: "Test Title",
				Body:  "Test Body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.idParam)

			mockClient := &mockPostGetter{
				getPostFunc: func(ctx context.Context, postID int) (*client.Resp, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			h := &Handler{
				Client: mockClient,
				Logger: slogdiscard.NewDiscardLogger(),
				Config: &config.Config{HTTPClient: config.HTTPClientConfig{Timeout: time.Second * 5}},
			}

			err := h.Proxy(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var actualBody any

			switch tt.expectedBody.(type) {
			case ErrorResponse:
				var got ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &got)
				require.NoError(t, err, "failed to unmarshal error_response body")
				actualBody = got
			case client.Resp:
				var got client.Resp
				err = json.Unmarshal(rec.Body.Bytes(), &got)
				require.NoError(t, err, "failed to unmarshal response body")
				actualBody = got
			}
			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}
