package amidware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/jpfluger/alibs-slim/asessions"
)

// IActionRoutesProvisioner defines the interface for route provisioners.
type IActionRoutesProvisioner interface {
	MatchesWhitelistActionPath(targetPath string) bool
	LogError(c echo.Context, err error)
	GetUrlNotFoundActions() string
	IsMaintenanceMode() bool
	GetUrlMaintenance() string
}

// ActionRoutesConfig holds the configuration for action routes middleware.
type ActionRoutesConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
	// ActionKeyUrls holds the mapping from action keys to URLs.
	ActionKeyUrls asessions.ActionKeyUrls
	// Provisioner is the interface implementation for route provisioning.
	Provisioner IActionRoutesProvisioner
}

// NewActionRoutes creates a new ActionRoutes middleware with the provided configuration.
func NewActionRoutes(akUrls asessions.ActionKeyUrls, provisioner IActionRoutesProvisioner) echo.MiddlewareFunc {
	// Initialize the ActionRoutesConfig with default values and provided arguments.
	config := &ActionRoutesConfig{
		Skipper:       middleware.DefaultSkipper,
		ActionKeyUrls: akUrls,
		Provisioner:   provisioner,
	}

	// Validate the configuration.
	if err := validateActionRoutesConfig(config); err != nil {
		panic(err)
	}

	return ActionRoutes(config)
}

// validateActionRoutesConfig checks the provided configuration for any missing or invalid values.
func validateActionRoutesConfig(config *ActionRoutesConfig) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if config.Provisioner == nil {
		return fmt.Errorf("Provisioner is nil")
	}
	if config.ActionKeyUrls == nil || len(config.ActionKeyUrls) == 0 {
		return fmt.Errorf("ActionKeyUrls is empty")
	}
	for _, aku := range config.ActionKeyUrls {
		if aku.Key.IsEmpty() {
			return fmt.Errorf("ActionKeyUrls.Key is empty")
		}
		if strings.TrimSpace(aku.Url) == "" {
			return fmt.Errorf("ActionKeyUrls.Url is empty")
		}
	}
	return nil
}

// ActionRoutes returns a middleware function that handles action routing based on user session status.
func ActionRoutes(config *ActionRoutesConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Maintenance mode check using Provisioner
			if config.Provisioner.IsMaintenanceMode() {
				maintenanceUrl := config.Provisioner.GetUrlMaintenance()
				if c.Request().URL.Path != maintenanceUrl {
					return c.Redirect(http.StatusFound, maintenanceUrl)
				}
			}

			us := asessions.CastLoginSessionPermFromEchoContext(c)
			if us == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "actions: user session not found")
			}

			if us.GetStatusType() == asessions.LOGIN_SESSION_STATUS_ACTIONS {
				mainUrl, checkPrefix := determineActionURL(us, config)

				if !isActionURLValid(c, mainUrl, checkPrefix, config) {
					return c.Redirect(http.StatusFound, mainUrl)
				}
			}

			return next(c)
		}
	}
}

// determineActionURL determines the main URL and whether to check the prefix based on the user session.
func determineActionURL(us asessions.ILoginSessionPerm, config *ActionRoutesConfig) (string, bool) {
	mainUrl := config.Provisioner.GetUrlNotFoundActions()
	checkPrefix := false

	if len(us.GetActions()) > 0 {
		aku := config.ActionKeyUrls.Find(us.GetActions()[0])
		if aku != nil {
			mainUrl = aku.Url
			checkPrefix = aku.CheckPrefix
		}
	}

	return mainUrl, checkPrefix
}

// isActionURLValid checks if the requested URL matches the action URL or is whitelisted.
func isActionURLValid(c echo.Context, mainUrl string, checkPrefix bool, config *ActionRoutesConfig) bool {
	if checkPrefix && strings.HasPrefix(c.Request().URL.Path, mainUrl) {
		return true
	}

	if mainUrl == c.Request().URL.Path {
		return true
	}

	return config.Provisioner.MatchesWhitelistActionPath(c.Request().URL.Path)
}
