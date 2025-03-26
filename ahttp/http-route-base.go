package ahttp

import (
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/asessions"
	"sync"
)

// IRoute defines the interface for route-related information.
// It is implemented by types that represent specific routes within the application.
type IRoute interface {
	GetRouteId() HttpRouteId         // Returns the unique identifier of the route.
	GetMethod() HttpMethod           // Returns the HTTP method (GET, POST, etc.) associated with the route.
	GetPath() string                 // Returns the URL path of the route.
	GetPerms() asessions.PermSet     // Returns the permissions required to access the route.
	GetRouteNotFoundId() HttpRouteId // Returns the identifier of the route to use when the current route is not found.
	GetIndexAdded() int
	CreateHandler() echo.HandlerFunc
	GetWhitelist() string
}

// IRoutes is a slice of IRoute interfaces.
type IRoutes []IRoute

// ToMap converts a slice of IRoute interfaces into a map keyed by the route ID.
func (rts IRoutes) ToMap() IRouteMap {
	rtMap := IRouteMap{}
	if rts == nil || len(rts) == 0 {
		return rtMap
	}
	for _, route := range rts {
		rtMap[route.GetRouteId()] = route
	}
	return rtMap
}

// IRouteMap is a map that associates route IDs with their corresponding IRoute interface.
type IRouteMap map[HttpRouteId]IRoute

// RouteBase provides a base implementation of the IRoute interface.
type RouteBase struct {
	RouteId         HttpRouteId       // Unique identifier of the route.
	Method          HttpMethod        // HTTP method associated with the route.
	Path            string            // URL path of the route.
	Perms           asessions.PermSet // Permissions required to access the route.
	RouteNotFoundId HttpRouteId       // Identifier of the route to use when the current route is not found.

	Whitelist string // whitelist is needed when using actions or maintenance mode.

	indexAdded int          // The order in which the route was added.
	mu         sync.RWMutex // Mutex for safe concurrent access to the route's fields.
}

// GetIndexAdded returns the position at which the route was added.
func (rb *RouteBase) GetIndexAdded() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.indexAdded
}

// GetRouteId returns the route's unique identifier.
func (rb *RouteBase) GetRouteId() HttpRouteId {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.RouteId
}

// GetMethod returns the route's HTTP method.
func (rb *RouteBase) GetMethod() HttpMethod {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.Method
}

// GetPath returns the route's URL path.
func (rb *RouteBase) GetPath() string {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.Path
}

// GetPerms returns the permissions required to access the route.
func (rb *RouteBase) GetPerms() asessions.PermSet {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.Perms
}

// GetRouteNotFoundId returns the identifier of the route to use when the current route is not found.
func (rb *RouteBase) GetRouteNotFoundId() HttpRouteId {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.RouteNotFoundId
}

// GetWhitelist returns the whitelisted path.
func (rb *RouteBase) GetWhitelist() string {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.Whitelist
}
