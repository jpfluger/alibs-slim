package ahttp

import (
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/asessions"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRoute is a mock implementation of the IRoute interface for testing.
type MockRoute struct {
	routeId         HttpRouteId
	method          HttpMethod
	path            string
	perms           []string
	routeNotFoundId HttpRouteId
	indexAdded      int
	whitelist       string
}

func (m *MockRoute) GetRouteId() HttpRouteId { return m.routeId }
func (m *MockRoute) GetMethod() HttpMethod   { return m.method }
func (m *MockRoute) GetPath() string         { return m.path }
func (m *MockRoute) GetPerms() asessions.PermSet {
	var mockPSet asessions.PermSet
	if m.perms != nil && len(m.perms) > 0 {
		pset := asessions.PermSet{}
		for _, perm := range m.perms {
			pset.SetPerm(asessions.NewPerm(perm))
		}
		pset = pset
	}
	return mockPSet
}
func (m *MockRoute) GetRouteNotFoundId() HttpRouteId { return m.routeNotFoundId }
func (m *MockRoute) GetIndexAdded() int              { return m.indexAdded }
func (m *MockRoute) CreateHandler() echo.HandlerFunc {
	return func(c echo.Context) error { return nil }
}
func (m *MockRoute) GetWhitelist() string { return m.whitelist }

func setup() *HttpRouteMap {
	return &HttpRouteMap{
		RWMutex: sync.RWMutex{},
		routes:  make(map[HttpRouteId]IRoute),
		urlHome: "/home",
		urlDash: "/dashboard",
	}
}

func TestHttpRouteMap_Get(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "home",
		path:    "/home",
	}
	m.Set("home", route)

	assert.Equal(t, route, m.Get("home"))
	assert.Nil(t, m.Get("nonexistent"))
}

func TestHttpRouteMap_Set(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "home",
		path:    "/home",
	}
	m.Set("home", route)

	assert.Equal(t, route, m.Get("home"))
}

func TestHttpRouteMap_MustUrl(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "home",
		path:    "/home",
	}
	m.Set("home", route)

	assert.Equal(t, "/home", m.MustUrl("home"))
	assert.Equal(t, "/home", m.MustUrl("nonexistent"))
}

func TestHttpRouteMap_MustUrlElseDash(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "dashboard",
		path:    "/dashboard",
	}
	m.Set("dashboard", route)

	assert.Equal(t, "/dashboard", m.MustUrlElseDash("dashboard"))
	assert.Equal(t, "/dashboard", m.MustUrlElseDash("nonexistent"))
}

func TestHttpRouteMap_MustUrlRedirect(t *testing.T) {
	m := setup()

	assert.Equal(t, "/home", m.MustUrlRedirect(false))
	assert.Equal(t, "/dashboard", m.MustUrlRedirect(true))
}

func TestHttpRouteMap_SetRouteMapCheckKey(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "home",
		path:    "/home",
	}
	err := m.Set("home", route)
	assert.NoError(t, err)

	err = m.Set("home", route)
	assert.Error(t, err)
}

func TestHttpRouteMap_SetRouteMapSpecials(t *testing.T) {
	m := setup()
	routeHome := &MockRoute{
		routeId: "home",
		path:    "/home",
	}
	routeDash := &MockRoute{
		routeId: "dashboard",
		path:    "/dashboard",
	}
	routeNoLogin := &MockRoute{
		routeId: "noLogin",
		path:    "/noLogin",
	}
	routeInvalidPerms := &MockRoute{
		routeId: "invalidPerms",
		path:    "/invalidPerms",
	}
	m.Set("home", routeHome)
	m.Set("dashboard", routeDash)
	m.Set("noLogin", routeNoLogin)
	m.Set("invalidPerms", routeInvalidPerms)

	err := m.SetRouteMapSpecials("home", "dashboard", "noLogin", "invalidPerms")
	assert.NoError(t, err)

	err = m.SetRouteMapSpecials("nonexistent", "dashboard", "noLogin", "invalidPerms")
	assert.Error(t, err)
}

func TestHttpRouteMap_GetUrlNoLogin(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "noLogin",
		path:    "/noLogin",
	}
	m.Set("noLogin", route)
	err := m.SetRouteMapSpecials("home", "dashboard", "noLogin", "invalidPerms")
	assert.Error(t, err)
	assert.Equal(t, "/noLogin", m.GetUrlNoLogin())
}

func TestHttpRouteMap_GetUrlInvalidPerms(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "invalidPerms",
		path:    "/invalidPerms",
	}
	m.Set("invalidPerms", route)
	err := m.SetRouteMapSpecials("home", "dashboard", "noLogin", "invalidPerms")
	assert.Error(t, err)
	assert.Equal(t, "/invalidPerms", m.GetUrlInvalidPerms())
}

func TestHttpRouteMap_DeleteRoute(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "home",
		path:    "/home",
	}
	m.Set("home", route)

	err := m.DeleteRoute("home")
	assert.NoError(t, err)
	assert.Nil(t, m.Get("home"))

	err = m.DeleteRoute("nonexistent")
	assert.Error(t, err)
}

func TestHttpRouteMap_RouteExists(t *testing.T) {
	m := setup()
	route := &MockRoute{
		routeId: "home",
		path:    "/home",
	}
	m.Set("home", route)

	assert.True(t, m.RouteExists("home"))
	assert.False(t, m.RouteExists("nonexistent"))
}
