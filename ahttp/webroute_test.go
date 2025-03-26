package ahttp

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/asessions"
	"testing"
)

// MockRouteBase is a mock implementation of RouteBase for testing purposes.
type MockRouteBase struct {
	RouteId HttpRouteId
	Method  HttpMethod
	Path    string
	Perms   asessions.PermSet
}

// MockHandler is a simple echo.HandlerFunc for testing purposes.
func MockHandler(echo.Context) error {
	return nil
}

// TestNewWebRouteSPerm tests the NewWebRouteSPerm function for creating a new WebRoute with string permissions.
func TestNewWebRouteSPerm(t *testing.T) {
	mockRouteId := HttpRouteId("testRouteId")
	mockMethod := HTTPMETHOD_GET
	mockURL := "/test"
	mockPerms := []string{"read:R", "write:CU"}
	mockCreateRouteHandler := func(route IRoute) echo.HandlerFunc {
		return MockHandler
	}

	//webRoute := NewWebRouteSPerm(mockRouteId, mockMethod, mockURL, mockPerms, mockCreateRouteHandler)
	webRoute := NewWRPermStr(mockRouteId, mockMethod, mockURL, mockPerms, mockCreateRouteHandler)

	assert.NotNil(t, webRoute, "WebRoute should not be nil")
	assert.Equal(t, mockRouteId, webRoute.RouteId, "RouteId should match")
	assert.Equal(t, mockMethod, webRoute.Method, "Method should match")
	assert.Equal(t, mockURL, webRoute.Path, "Path should match")
	assert.NotNil(t, webRoute.Perms, "Perms should not be nil")
	assert.Equal(t, 2, len(webRoute.Perms), "Perms should have two permissions")
	assert.NotNil(t, webRoute.handler, "Handler should not be nil")
}

// TestNewWebRoute tests the NewWebRoute function for creating a new WebRoute with PermSet permissions.
func TestNewWebRoute(t *testing.T) {
	mockRouteId := HttpRouteId("testRouteId")
	mockMethod := HTTPMETHOD_POST
	mockURL := "/test"
	mockPerms := []string{"execute:X"}
	mockCreateRouteHandler := func(route IRoute) echo.HandlerFunc {
		return MockHandler
	}

	webRoute := NewWRPermStr(mockRouteId, mockMethod, mockURL, mockPerms, mockCreateRouteHandler)

	assert.NotNil(t, webRoute, "WebRoute should not be nil")
	assert.Equal(t, mockRouteId, webRoute.RouteId, "RouteId should match")
	assert.Equal(t, mockMethod, webRoute.Method, "Method should match")
	assert.Equal(t, mockURL, webRoute.Path, "Path should match")
	assert.NotNil(t, webRoute.Perms, "Perms should not be nil")
	assert.Equal(t, 1, len(webRoute.Perms), "Perms should have one permission")
	assert.NotNil(t, webRoute.handler, "Handler should not be nil")
}
