package ahttp

// HttpRouteIdMap defines a map that associates HttpRouteId keys with URL strings.
type HttpRouteIdMap map[HttpRouteId]string

// NewRouteRedirectDefaults creates a new HttpRouteIdMap with default route redirects.
func NewRouteRedirectDefaults() HttpRouteIdMap {
	return HttpRouteIdMap{
		ROUTE_REDIRECT_DEFAULT:    "/", // Default redirect route (typically to home page).
		ROUTE_REDIRECT_AUTH_ERROR: "/", // Redirect route for authentication errors.
	}
}
