package ahttp

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

//// MockRoute is a mock implementation of the IRoute interface for testing.
//type MockRoute struct {
//	routeId         HttpRouteId
//	method          HttpMethod
//	path            string
//	perms           []string
//	routeNotFoundId HttpRouteId
//	indexAdded      int
//}
//
//func (m *MockRoute) GetRouteId() HttpRouteId         { return m.routeId }
//func (m *MockRoute) GetMethod() HttpMethod           { return m.method }
//func (m *MockRoute) GetPath() string                 { return m.path }
//func (m *MockRoute) GetPerms() []string              { return m.perms }
//func (m *MockRoute) GetRouteNotFoundId() HttpRouteId { return m.routeNotFoundId }
//func (m *MockRoute) GetIndexAdded() int              { return m.indexAdded }
//func (m *MockRoute) CreateHandler() echo.HandlerFunc {
//	return func(c echo.Context) error { return nil }
//}

func setupWebRouteManager() *WebRouteManager {
	return &WebRouteManager{
		HttpRouteMap: HttpRouteMap{
			RWMutex: sync.RWMutex{},
			routes:  make(map[HttpRouteId]IRoute),
		},
		allowedActionPaths:      []string{},
		authenticateProvisioner: nil,
	}
}

func TestWebRouteManager_AddRoute(t *testing.T) {
	wrm := setupWebRouteManager()
	route := &MockRoute{
		routeId:    "home",
		method:     HTTPMETHOD_GET,
		path:       "/home",
		indexAdded: 1,
	}

	err := wrm.AddRoute(route)
	assert.NoError(t, err)
	assert.Equal(t, route, wrm.Get("home"))
}

func TestWebRouteManager_AddRoute_Duplicate(t *testing.T) {
	wrm := setupWebRouteManager()
	route1 := &MockRoute{
		routeId:    "home",
		method:     HTTPMETHOD_GET,
		path:       "/home",
		indexAdded: 1,
	}
	route2 := &MockRoute{
		routeId:    "home",
		method:     HTTPMETHOD_GET,
		path:       "/home",
		indexAdded: 2,
	}

	err := wrm.AddRoute(route1)
	assert.NoError(t, err)

	err = wrm.AddRoute(route2)
	assert.Error(t, err)
}

func TestWebRouteManager_InitRoutesWithEcho(t *testing.T) {
	e := echo.New()
	wrm := setupWebRouteManager()
	route := &MockRoute{
		routeId:    "home",
		method:     HTTPMETHOD_GET,
		path:       "/home",
		indexAdded: 1,
	}

	err := wrm.AddRoute(route)
	assert.NoError(t, err)

	err = wrm.InitRoutesWithEcho(e)
	assert.NoError(t, err)

	// Check if the route is registered in Echo
	found := false
	for _, r := range e.Routes() {
		if r.Path == "/home" && r.Method == "GET" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestWebRouteManager_AddWhitelistActionPath(t *testing.T) {
	wrm := setupWebRouteManager()
	wrm.AddWhitelistActionPath("/whitelist")

	assert.Contains(t, wrm.allowedActionPaths, "/whitelist")
}

func TestWebRouteManager_MatchesWhitelistActionPath(t *testing.T) {
	wrm := setupWebRouteManager()
	wrm.AddWhitelistActionPath("/whitelist")

	assert.True(t, wrm.MatchesWhitelistActionPath("/whitelist"))
	assert.False(t, wrm.MatchesWhitelistActionPath("/notwhitelist"))
}

func TestWebRouteManager_LogError(t *testing.T) {
	wrm := setupWebRouteManager()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := fmt.Errorf("test error")
	wrm.LogError(c, err)
	// Check logs manually or use a logging library that supports testing
}

func TestWebRouteManager_LogAuthError(t *testing.T) {
	wrm := setupWebRouteManager()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := fmt.Errorf("auth error")
	wrm.LogAuthError(c, err)
	// Check logs manually or use a logging library that supports testing
}
