package ahttp

import (
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/asessions"
	"strings"
)

// CreateRouteHandler defines a function type that takes an IRoute and returns an echo.HandlerFunc.
// This allows for flexible handling of route creation with custom logic.
type CreateRouteHandler func(route IRoute) echo.HandlerFunc

// CreateRouteHandlerByEchoHandlerFunc wraps a given echo.HandlerFunc into a CreateRouteHandler type.
// It panics if the handler is nil, ensuring a valid handler is always provided.
func CreateRouteHandlerByEchoHandlerFunc(handler echo.HandlerFunc) CreateRouteHandler {
	if handler == nil {
		panic("handler is nil")
	}
	return func(route IRoute) echo.HandlerFunc {
		return handler
	}
}

// WebRoute is a struct that implements IRoute using an echo.HandlerFunc.
// It represents a web route with a handler and base routing information.
type WebRoute struct {
	RouteBase
	handler echo.HandlerFunc
}

// CreateHandler returns the handler function associated with the WebRoute.
// This function is part of the WebRoute's implementation of the IRoute interface.
func (wr *WebRoute) CreateHandler() echo.HandlerFunc {
	return wr.handler
}

// NewWRPermSetEH creates a new WebRoute with permission sets based on the specified handler.
func NewWRPermSetEH(httpRouteId HttpRouteId, method HttpMethod, url string, permSet asessions.PermSet, handler echo.HandlerFunc) *WebRoute {
	return NewWRPermSet(httpRouteId, method, url, permSet, CreateRouteHandlerByEchoHandlerFunc(handler))
}

// NewWRPermSetEHAutoWL creates a new WebRoute with permission sets and automatic whitelisting.
func NewWRPermSetEHAutoWL(httpRouteId HttpRouteId, method HttpMethod, url string, permSet asessions.PermSet, handler echo.HandlerFunc) *WebRoute {
	return NewWRPermSetAutoWL(httpRouteId, method, url, permSet, CreateRouteHandlerByEchoHandlerFunc(handler))
}

// NewWRPermSetEHWL creates a new WebRoute with permission sets and a specified whitelist.
func NewWRPermSetEHWL(httpRouteId HttpRouteId, method HttpMethod, url string, permSet asessions.PermSet, handler echo.HandlerFunc, whitelist string) *WebRoute {
	return NewWRPermSetWL(httpRouteId, method, url, permSet, CreateRouteHandlerByEchoHandlerFunc(handler), whitelist)
}

// NewWRPermStrEH creates a new WebRoute with permissions defined by strings and a handler.
func NewWRPermStrEH(httpRouteId HttpRouteId, method HttpMethod, url string, permStr []string, handler echo.HandlerFunc) *WebRoute {
	return NewWRPermStr(httpRouteId, method, url, permStr, CreateRouteHandlerByEchoHandlerFunc(handler))
}

// NewWRPermStrEHAutoWL creates a new WebRoute with string permissions and automatic whitelisting.
func NewWRPermStrEHAutoWL(httpRouteId HttpRouteId, method HttpMethod, url string, permStr []string, handler echo.HandlerFunc) *WebRoute {
	return NewWRPermStrAutoWL(httpRouteId, method, url, permStr, CreateRouteHandlerByEchoHandlerFunc(handler))
}

// NewWRPermStrEHWL creates a new WebRoute with string permissions and a specified whitelist.
func NewWRPermStrEHWL(httpRouteId HttpRouteId, method HttpMethod, url string, permStr []string, handler echo.HandlerFunc, whitelist string) *WebRoute {
	return NewWRPermStrWL(httpRouteId, method, url, permStr, CreateRouteHandlerByEchoHandlerFunc(handler), whitelist)
}

// NewWRPermStr initializes a WebRoute with string permissions and a route handler.
func NewWRPermStr(httpRouteId HttpRouteId, method HttpMethod, url string, permStr []string, createRouteHandler CreateRouteHandler) *WebRoute {
	return NewWROptions(httpRouteId, method, url, nil, permStr, createRouteHandler, false, "")
}

// NewWRPermStrAutoWL initializes a WebRoute with string permissions and automatic whitelisting.
func NewWRPermStrAutoWL(httpRouteId HttpRouteId, method HttpMethod, url string, permStr []string, createRouteHandler CreateRouteHandler) *WebRoute {
	return NewWROptions(httpRouteId, method, url, nil, permStr, createRouteHandler, true, "")
}

// NewWRPermStrWL initializes a WebRoute with string permissions and a specific whitelist.
func NewWRPermStrWL(httpRouteId HttpRouteId, method HttpMethod, url string, permStr []string, createRouteHandler CreateRouteHandler, whitelist string) *WebRoute {
	return NewWROptions(httpRouteId, method, url, nil, permStr, createRouteHandler, false, whitelist)
}

// NewWRPermSet initializes a WebRoute with a permission set and route handler.
func NewWRPermSet(httpRouteId HttpRouteId, method HttpMethod, url string, permSet asessions.PermSet, createRouteHandler CreateRouteHandler) *WebRoute {
	return NewWROptions(httpRouteId, method, url, permSet, nil, createRouteHandler, false, "")
}

// NewWRPermSetAutoWL initializes a WebRoute with a permission set and automatic whitelisting.
func NewWRPermSetAutoWL(httpRouteId HttpRouteId, method HttpMethod, url string, permSet asessions.PermSet, createRouteHandler CreateRouteHandler) *WebRoute {
	return NewWROptions(httpRouteId, method, url, permSet, nil, createRouteHandler, true, "")
}

// NewWRPermSetWL initializes a WebRoute with a permission set and a specific whitelist.
func NewWRPermSetWL(httpRouteId HttpRouteId, method HttpMethod, url string, permSet asessions.PermSet, createRouteHandler CreateRouteHandler, whitelist string) *WebRoute {
	return NewWROptions(httpRouteId, method, url, permSet, nil, createRouteHandler, false, whitelist)
}

// NewWROptions creates a new WebRoute with specified options for permissions, handlers, and whitelisting.
func NewWROptions(httpRouteId HttpRouteId, method HttpMethod, url string, permSet asessions.PermSet, permStr []string, createRouteHandler CreateRouteHandler, doAutoWhitelist bool, whitelist string) *WebRoute {
	// Convert permissions from strings to PermSet if provided as a list of strings.
	if permSet == nil || len(permSet) == 0 {
		if permStr != nil || len(permStr) > 0 {
			permSet = asessions.NewPermSetByString(permStr)
		}
	}
	// Trim and set the whitelist, using the URL if auto-whitelisting is enabled.
	whitelist = strings.TrimSpace(whitelist)
	if doAutoWhitelist && whitelist != "" {
		whitelist = url
	}

	// Initialize a WebRoute instance with the provided options.
	wr := &WebRoute{
		RouteBase: RouteBase{
			RouteId:   httpRouteId,
			Method:    method,
			Path:      url,
			Perms:     nil,
			Whitelist: strings.TrimSpace(whitelist),
		},
	}

	// Assign permissions if provided, defaulting to an empty PermSet if none is provided.
	if permSet != nil {
		wr.Perms = permSet
	} else {
		wr.Perms = asessions.PermSet{}
	}

	// Ensure a route handler is provided, setting it to the WebRoute if valid.
	if createRouteHandler == nil {
		panic("createRouteHandler is nil")
	} else {
		wr.handler = createRouteHandler(wr)
	}

	return wr
}
