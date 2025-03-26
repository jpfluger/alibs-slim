package ahttp

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/alog"
	"github.com/jpfluger/alibs-slim/amidware"
	"strings"
	"sync"
)

// IWebRouteManager defines the interface for managing web routes.
type IWebRouteManager interface {
	GetAuthenticateProvisioner() amidware.IAuthenticateProvisioner
	SetAuthenticateProvisioner(provisioner amidware.IAuthenticateProvisioner)
	AddRoute(e *echo.Echo, route IRoute) error
	SetRouteSpecials(rHome HttpRouteId, rDash HttpRouteId, rNoLogin HttpRouteId, rInvalidPerms HttpRouteId) error
	InitRoutesWithEcho(e *echo.Echo) error
	LogError(c echo.Context, err error)
	LogAuthError(c echo.Context, err error)
}

// IWebRouteController defines the interface for WebRouteManager in production context
type IWebRouteController interface {
	Get(key HttpRouteId) IRoute
	MustUrl(key HttpRouteId) string
	MustUrlElseDash(key HttpRouteId) string
	MustUrlRedirect(isLoggedIn bool) string
	MustUrlElseIsLogin(key HttpRouteId, isLoggedIn bool) string
	GetUrlNoLogin() string
	GetUrlInvalidPerms() string
	RouteExists(key HttpRouteId) bool
	LogError(c echo.Context, err error)
	LogAuthError(c echo.Context, err error)
	MatchesWhitelistActionPath(targetPath string) bool
}

// WebRouteManager manages the routing for web requests.
type WebRouteManager struct {
	HttpRouteMap                                              // Embedding HttpRouteMap for route management.
	allowedActionPaths      []string                          // List of paths that are whitelisted.
	authenticateProvisioner amidware.IAuthenticateProvisioner // Provisioner for authentication.
	mu                      sync.RWMutex                      // Mutex for concurrent access control.
}

func NewWebRouteManager() *WebRouteManager {
	return &WebRouteManager{
		HttpRouteMap: HttpRouteMap{
			RWMutex: sync.RWMutex{},
			routes:  map[HttpRouteId]IRoute{},
		},
		allowedActionPaths:      []string{},
		authenticateProvisioner: nil,
		mu:                      sync.RWMutex{},
	}
}

// GetAuthenticateProvisioner retrieves the provisioner for authentication.
func (wrm *WebRouteManager) GetAuthenticateProvisioner() amidware.IAuthenticateProvisioner {
	return wrm.authenticateProvisioner
}

// SetAuthenticateProvisioner sets the provisioner for authentication.
func (wrm *WebRouteManager) SetAuthenticateProvisioner(provisioner amidware.IAuthenticateProvisioner) {
	wrm.authenticateProvisioner = provisioner
}

// AddRoute adds a new route to the manager's route map.
func (wrm *WebRouteManager) AddRoute(route IRoute) error {
	// Validate route.
	if route == nil || route.GetRouteId().IsEmpty() {
		return fmt.Errorf("route is nil or routeid is empty")
	}

	// Trim spaces from the URL path.
	url := strings.TrimSpace(route.GetPath())
	if url == "" {
		return fmt.Errorf("url is empty")
	}

	if route.GetMethod().IsEmpty() {
		return fmt.Errorf("method is empty for route '%s'", route.GetRouteId().String())
	}

	// Check for duplicate routes.
	compare := wrm.Get(route.GetRouteId())
	if compare != nil {
		if compare.GetMethod() == route.GetMethod() {
			return fmt.Errorf("duplicate route '%s' detected for same method '%s'", route.GetRouteId().String(), compare.GetMethod().String())
		}
	}

	// Add the route to the map.
	if err := wrm.Set(route.GetRouteId(), route); err != nil {
		return err
	}

	return nil
}

// InitRoutesWithEcho initializes routes with the Echo instance.
func (wrm *WebRouteManager) InitRoutesWithEcho(e *echo.Echo) error {
	if e == nil {
		return fmt.Errorf("echo instance (e) is nil")
	}
	arr := wrm.GetRoutesArray()
	if len(arr) == 0 {
		return fmt.Errorf("no routes available")
	}
	for _, route := range arr {
		// Create the handler for the route.
		handler := route.CreateHandler()
		if handler == nil {
			return fmt.Errorf("handler is nil for route '%s'", route.GetRouteId().String())
		}
		// Set up authentication middleware if permissions are provided.
		var mwAuth echo.MiddlewareFunc
		if perms := route.GetPerms(); perms != nil && len(perms) > 0 {
			mwAuth = amidware.NewAuthenticatePermConfig(perms, wrm.authenticateProvisioner)
		}

		// Prepare middleware slice if mwAuth is not nil
		middlewares := []echo.MiddlewareFunc{}
		if mwAuth != nil {
			middlewares = append(middlewares, mwAuth)
		}

		// Trim spaces from the URL path.
		url := strings.TrimSpace(route.GetPath())

		// Register the route with the appropriate HTTP method.
		var eroute *echo.Route
		switch route.GetMethod() {
		case HTTPMETHOD_GET:
			eroute = e.GET(url, handler, middlewares...)
		case HTTPMETHOD_POST:
			eroute = e.POST(url, handler, middlewares...)
		case HTTPMETHOD_PUT:
			eroute = e.PUT(url, handler, middlewares...)
		case HTTPMETHOD_DELETE:
			eroute = e.DELETE(url, handler, middlewares...)
		case HTTPMETHOD_PATCH:
			eroute = e.PATCH(url, handler, middlewares...)
		default:
			return fmt.Errorf("unknown http method '%s' for route '%s'", route.GetMethod().String(), route.GetRouteId().String())
		}

		// Set a custom name for the route.
		eroute.Name = fmt.Sprintf("%s_%s", route.GetRouteId(), route.GetMethod().String())

		// Try adding to the whitelist. If empty, nothing is added.
		wrm.addWhitelistActionPath(route.GetWhitelist())
	}

	return nil
}

// AddWhitelistActionPath adds a path to the whitelist, applying a mutex lock.
func (wrm *WebRouteManager) AddWhitelistActionPath(target string) {
	wrm.mu.Lock()
	defer wrm.mu.Unlock()
	wrm.addWhitelistActionPath(target)
}

// addWhitelistActionPath adds a path to the whitelist.
func (wrm *WebRouteManager) addWhitelistActionPath(target string) {
	// Don't add empty targets.
	if strings.TrimSpace(target) == "" {
		return
	}
	// Check if the path is already in the whitelist.
	for _, pth := range wrm.allowedActionPaths {
		if pth == target {
			return
		}
	}
	wrm.allowedActionPaths = append(wrm.allowedActionPaths, target)
}

// MatchesWhitelistActionPath checks if a given path matches any in the whitelist.
func (wrm *WebRouteManager) MatchesWhitelistActionPath(targetPath string) bool {
	for _, pth := range wrm.allowedActionPaths {
		if strings.HasPrefix(targetPath, pth) {
			return true
		}
	}
	return false
}

// LogError logs an error.
func (wrm *WebRouteManager) LogError(c echo.Context, err error) {
	alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("wrm.LogError")
}

// LogAuthError logs an authentication error.
func (wrm *WebRouteManager) LogAuthError(c echo.Context, err error) {
	alog.LOGGER(alog.LOGGER_AUTH).Err(err).Msg("wrm.LogAuthError")
}
