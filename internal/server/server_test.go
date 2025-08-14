package server

import (
	"github.com/4uvirik/ProxyForPublicAPI/internal/handlers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupRouts(t *testing.T) {

	h := &handlers.Handler{}

	e := SetupRouts(h)

	assert.NotNil(t, e, "Echo instance should not be nil")

	routes := e.Routes()
	assert.NotEmpty(t, routes, "Routes should not be empty")

	// Проверка маршрутов в routes_tests.go
}
