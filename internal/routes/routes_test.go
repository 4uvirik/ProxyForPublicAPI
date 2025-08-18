package routes

import (
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitRoutes(t *testing.T) {
	h := &handlers.Handler{}
	e := echo.New()
	InitRoutes(e, h)

	routes := e.Routes()
	found := true

	for _, r := range routes {
		if r.Method == "GET" && r.Path == "/posts/:id" {
			found = true
			break
		}
	}

	assert.True(t, found, "GET request and /posts/:id route should be registered")

}
