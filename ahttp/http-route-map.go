package ahttp

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"sort"
	"sync"
)

// HttpRouteMap is a struct that associates HttpRouteId with an IRoute value.
type HttpRouteMap struct {
	sync.RWMutex
	routes          map[HttpRouteId]IRoute
	urlHome         string
	urlDash         string
	urlNoLogin      string
	urlInvalidPerms string
}

// Get retrieves the IRoute associated with the given HttpRouteId.
// It uses a read lock to allow concurrent read access.
func (m *HttpRouteMap) Get(key HttpRouteId) IRoute {
	m.RLock()         // Acquire a read lock.
	defer m.RUnlock() // Defer the unlock to the end of the function.
	return m.routes[key]
}

// Set sets the IRoute for a given HttpRouteId in the map, ensuring the path is not set by the key.
// It uses a write lock to prevent concurrent writes.
func (m *HttpRouteMap) Set(key HttpRouteId, value IRoute) error {
	m.Lock()         // Acquire a write lock.
	defer m.Unlock() // Defer the unlock to the end of the function.

	if _, ok := m.routes[key]; ok {
		return fmt.Errorf("duplicate route key %s", key)
	}
	m.routes[key] = value
	return nil
}

// MustUrl retrieves the URL associated with the given HttpRouteId or returns urlHome if not found.
func (m *HttpRouteMap) MustUrl(key HttpRouteId) string {
	m.RLock()
	defer m.RUnlock()
	return m.mustUrlElseIsLogin(key, false)
}

// MustUrlElseDash retrieves the URL associated with the given HttpRouteId or returns urlDash if not found.
func (m *HttpRouteMap) MustUrlElseDash(key HttpRouteId) string {
	m.RLock()
	defer m.RUnlock()
	return m.mustUrlElseIsLogin(key, true)
}

// MustUrlRedirect returns urlDash if the user is logged in, otherwise returns urlHome.
func (m *HttpRouteMap) MustUrlRedirect(isLoggedIn bool) string {
	m.RLock()
	defer m.RUnlock()
	return m.mustUrlElseIsLogin("", isLoggedIn)
}

// MustUrlElseIsLogin retrieves the URL associated with the given HttpRouteId or returns urlHome or urlDash based on login status.
func (m *HttpRouteMap) MustUrlElseIsLogin(key HttpRouteId, isLoggedIn bool) string {
	m.RLock()
	defer m.RUnlock()
	return m.mustUrlElseIsLogin(key, isLoggedIn)
}

// mustUrlElseIsLogin is an internal helper function to retrieve the URL or return a default based on login status.
func (m *HttpRouteMap) mustUrlElseIsLogin(key HttpRouteId, isLoggedIn bool) string {
	if key != "" {
		if route, exists := m.routes[key]; exists {
			if path := route.GetPath(); path != "" {
				return path
			}
		}
	}
	if isLoggedIn {
		return m.urlDash
	}
	return m.urlHome
}

// SetRouteMapCheckKey ensures the path is not set by the key.
// It uses a write lock to prevent concurrent writes and ensure atomicity.
//func (m *HttpRouteMap) SetRouteMapCheckKey(key HttpRouteId, route IRoute) error {
//	m.Lock()         // Acquire a write lock.
//	defer m.Unlock() // Defer the unlock to the end of the function.
//
//	if _, ok := m.routes[key]; ok {
//		return fmt.Errorf("duplicate route key %s", key)
//	}
//	m.routes[key] = route
//	return nil
//}

// SetRouteMapSpecials sets urlNoLogin and urlInvalidPerms.
// It uses a write lock to prevent concurrent writes and ensure atomicity.
func (m *HttpRouteMap) SetRouteMapSpecials(rHome HttpRouteId, rDash HttpRouteId, rNoLogin HttpRouteId, rInvalidPerms HttpRouteId) error {
	m.Lock()
	defer m.Unlock()
	var err error
	fnSetErr := func(myErr error) {
		if err == nil {
			err = myErr
		}
	}
	if route, exists := m.routes[rHome]; exists {
		m.urlHome = route.GetPath()
		if m.urlHome == "" {
			fnSetErr(fmt.Errorf("no route map for urlHome %s", rHome))
		}
	} else {
		fnSetErr(fmt.Errorf("no route map for urlHome %s", rHome))
	}
	if route, exists := m.routes[rDash]; exists {
		m.urlDash = route.GetPath()
		if m.urlDash == "" {
			fnSetErr(fmt.Errorf("no route map for urlDash %s", rDash))
		}
	} else {
		fnSetErr(fmt.Errorf("no route map for urlDash %s", rDash))
	}
	if route, exists := m.routes[rNoLogin]; exists {
		m.urlNoLogin = route.GetPath()
		if m.urlNoLogin == "" {
			fnSetErr(fmt.Errorf("no route map for urlNoLogin %s", rNoLogin))
		}
	} else {
		fnSetErr(fmt.Errorf("no route map for urlNoLogin %s", rNoLogin))
	}
	if route, exists := m.routes[rInvalidPerms]; exists {
		m.urlInvalidPerms = route.GetPath()
		if m.urlInvalidPerms == "" {
			fnSetErr(fmt.Errorf("no route map for urlInvalidPerms %s", rInvalidPerms))
		}
	} else {
		fnSetErr(fmt.Errorf("no route map for urlInvalidPerms %s", rInvalidPerms))
	}
	return err
}

// GetUrlNoLogin safely gets the urlNoLogin field from the global instance.
func (m *HttpRouteMap) GetUrlNoLogin() string {
	m.RLock()
	defer m.RUnlock()
	return m.urlNoLogin
}

// GetUrlInvalidPerms safely gets the urlInvalidPerms field from the global instance.
func (m *HttpRouteMap) GetUrlInvalidPerms() string {
	m.RLock()
	defer m.RUnlock()
	return m.urlInvalidPerms
}

// LogAuthError logs the authentication or authorization error.
func (m *HttpRouteMap) LogAuthError(c echo.Context, err error) {
	log.Err(err).Str("url", c.Request().URL.String()).Msg("failed to permit access")
}

// RouteExists checks if a route exists.
func (m *HttpRouteMap) RouteExists(key HttpRouteId) bool {
	m.RLock()
	defer m.RUnlock()
	_, exists := m.routes[key]
	return exists
}

// DeleteRoute deletes a route.
func (m *HttpRouteMap) DeleteRoute(key HttpRouteId) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	m.Lock()
	defer m.Unlock()

	if _, exists := m.routes[key]; !exists {
		return fmt.Errorf("key %s does not exist", key)
	}

	delete(m.routes, key)
	return nil
}

// GetRoutesArray returns an array of IRoute objects sorted by their index.
func (m *HttpRouteMap) GetRoutesArray() []IRoute {
	m.RLock()
	defer m.RUnlock()

	// Collect all IRoute objects from the map.
	routesArray := make([]IRoute, 0, len(m.routes))
	for _, route := range m.routes {
		routesArray = append(routesArray, route)
	}

	// Sort the array based on the GetIndexAdded() value.
	sort.Slice(routesArray, func(i, j int) bool {
		return routesArray[i].GetIndexAdded() < routesArray[j].GetIndexAdded()
	})

	return routesArray
}
