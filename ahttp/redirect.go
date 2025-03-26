package ahttp

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/aimage"
	"github.com/jpfluger/alibs-slim/arob"
	"github.com/jpfluger/alibs-slim/asessions"
	"net/http"
	"strings"
)

// ForceRedirectRawQuery redirects the client to the target URL, preserving the raw query parameters.
func ForceRedirectRawQuery(c echo.Context, targetUrl string) error {
	if targetUrl == "" {
		// If no target URL is provided, redirect to the root with the raw query.
		return c.Redirect(http.StatusFound, "/?"+c.Request().URL.RawQuery)
	}

	// Ensure the target URL ends with a slash.
	targetUrl = strings.TrimSuffix(targetUrl, "/") + "/"

	// Redirect to the target URL with the raw query appended.
	return c.Redirect(http.StatusFound, targetUrl+"?"+c.Request().URL.RawQuery)
}

// ForceRedirectRequestToBase64 redirects the client to the target URL with the request URL encoded in base64.
func ForceRedirectRequestToBase64(c echo.Context, targetUrl string) error {
	// Convert the request URL to a base64-encoded string.
	burl := aimage.ToBase64([]byte(c.Request().URL.String()))

	// Ensure the target URL ends with a slash.
	targetUrl = strings.TrimSuffix(targetUrl, "/") + "/"

	// Redirect to the target URL with the base64-encoded request URL as a parameter.
	return c.Redirect(http.StatusFound, targetUrl+"?burl="+burl)
}

// ResetUserSessionWithRedirect clears the user session and redirects to the given URL.
func ResetUserSessionWithRedirect(c echo.Context, manager *scs.SessionManager, urlTarget string) error {
	// Default to redirecting to the root if no target URL is provided.
	if urlTarget == "" {
		urlTarget = "/"
	}

	// Remove the user session and renew the session token.
	manager.Remove(c.Request().Context(), asessions.ECHOSCS_OBJECTKEY_USER_SESSION)
	_ = manager.RenewToken(c.Request().Context())

	// Respond with JSON if the request expects a JSON response, otherwise redirect.
	if IsRequestContentType(c, CHECK_MIME_TYPE_JSON) {
		return c.JSON(http.StatusOK, arob.NewROBWithRedirect(urlTarget))
	}
	return c.Redirect(http.StatusFound, urlTarget)
}

// ResetUserSessionWithRedirectJSONMessage clears the user session and returns a JSON message or redirects.
func ResetUserSessionWithRedirectJSONMessage(c echo.Context, manager *scs.SessionManager, urlTarget string, message arob.ROBMessage) error {
	// Default to redirecting to the root if no target URL is provided.
	if urlTarget == "" {
		urlTarget = "/"
	}

	// Remove the user session and renew the session token.
	manager.Remove(c.Request().Context(), asessions.ECHOSCS_OBJECTKEY_USER_SESSION)
	_ = manager.RenewToken(c.Request().Context())

	// Respond with JSON if the request expects a JSON response, otherwise redirect.
	if IsRequestContentType(c, CHECK_MIME_TYPE_JSON) {
		return c.JSON(http.StatusOK, arob.NewROBRedirectWithMessage(urlTarget, message))
	}
	return c.Redirect(http.StatusFound, urlTarget)
}
