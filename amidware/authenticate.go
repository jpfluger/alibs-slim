package amidware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/jpfluger/alibs-slim/asessions"
)

// IAuthenticateProvisioner defines the interface for authentication provisioning.
type IAuthenticateProvisioner interface {
	GetUrlNoLogin() string
	GetUrlInvalidPerms() string
	LogAuthError(c echo.Context, err error)
}

// AuthenticatePermConfig holds the configuration for permission-based authentication middleware.
type AuthenticatePermConfig struct {
	Skipper         middleware.Skipper       // Function to skip middleware.
	Perms           asessions.PermSet        // Set of permissions required for access.
	Provisioner     IAuthenticateProvisioner // Interface for provisioning URLs and logging.
	UrlNoLogin      string                   // URL to redirect to when no login is detected.
	UrlInvalidPerms string                   // URL to redirect to when permissions are invalid.
}

// NewAuthenticatePermConfig creates a new instance of AuthenticatePermConfig with default values.
func NewAuthenticatePermConfig(perms asessions.PermSet, provisioner IAuthenticateProvisioner) echo.MiddlewareFunc {
	return authenticateConfig(&AuthenticatePermConfig{
		Skipper:     middleware.DefaultSkipper,
		Perms:       perms,
		Provisioner: provisioner,
	})
}

// authenticateConfig returns a middleware function that enforces permission-based authentication.
func authenticateConfig(config *AuthenticatePermConfig) echo.MiddlewareFunc {
	if config == nil {
		panic("authenticateConfig: config cannot be nil")
	}

	if len(config.Perms) == 0 {
		panic("authenticateConfig: perms cannot be empty")
	}

	// Set URLs from the provisioner if not explicitly provided.
	if config.Provisioner != nil {
		if config.UrlNoLogin == "" {
			config.UrlNoLogin = config.Provisioner.GetUrlNoLogin()
		}
		if config.UrlInvalidPerms == "" {
			config.UrlInvalidPerms = config.Provisioner.GetUrlInvalidPerms()
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			us := asessions.CastLoginSessionPermFromEchoContext(c)
			if us == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "authenticator: session not found")
			}

			if us.GetStatusType() == asessions.LOGIN_SESSION_STATUS_NONE {
				err := echo.NewHTTPError(http.StatusUnauthorized, "authenticator: session status is not logged-in")
				config.Provisioner.LogAuthError(c, err)
				return c.Redirect(http.StatusFound, config.UrlNoLogin)
			}

			if !config.Perms.HasPermSet(us.GetPerms()) {
				err := echo.NewHTTPError(http.StatusForbidden, "authenticator: unauthorized session permission")
				config.Provisioner.LogAuthError(c, err)
				return c.Redirect(http.StatusFound, config.UrlInvalidPerms)
			}

			return next(c)
		}
	}
}
