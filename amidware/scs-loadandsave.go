package amidware

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/jpfluger/alibs-slim/asessions"
	"github.com/jpfluger/alibs-slim/autils"
	"net/http"
	"time"
)

// Echo-SCS-Session
// Modified from https://github.com/spazzymoto/echo-scs-session
// MIT License, Copyright (c) 2021 Robert Edwards

// SessionConfig holds the configuration for session management middleware.
type SessionConfig struct {
	Skipper             middleware.Skipper  // Function to skip middleware.
	SessionManager      *scs.SessionManager // Session manager instance from SCS.
	DefaultLanguageType autils.LanguageType // Default language type for new sessions.
	IsOnRequireSession  bool                // Flag to indicate if session creation is required.
}

// DefaultSessionConfig provides default settings for session management.
var DefaultSessionConfig = SessionConfig{
	Skipper: middleware.DefaultSkipper,
}

// SCSLoadAndSave initializes session management middleware with default configuration.
func SCSLoadAndSave(sessionManager *scs.SessionManager, isOnRequireSession bool) echo.MiddlewareFunc {
	c := DefaultSessionConfig
	c.SessionManager = sessionManager
	c.IsOnRequireSession = isOnRequireSession
	return SCSLoadAndSaveWithConfig(c)
}

// SCSLoadAndSaveWithConfig returns a middleware function that loads and saves session data.
func SCSLoadAndSaveWithConfig(config SessionConfig) echo.MiddlewareFunc {
	if config.SessionManager == nil {
		panic("SCSLoadAndSaveWithConfig: SessionManager cannot be nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			ctx := c.Request().Context()

			// If a new session, token is empty and err != nil.
			// If an existing session, token has value and err == nil.
			var token string
			cookie, err := c.Cookie(config.SessionManager.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err = config.SessionManager.Load(ctx, token)
			if err != nil {
				return err
			}

			c.SetRequest(c.Request().WithContext(ctx))

			if config.IsOnRequireSession {
				us, ok := config.SessionManager.Get(c.Request().Context(), asessions.ECHOSCS_OBJECTKEY_USER_SESSION).(asessions.UserSessionPerm)
				if !ok { // No? Create it.
					us = *asessions.NewUserSessionPerm()
					us.LanguageType = autils.GetHTTPAcceptedLanguageWithDefault(c.Request().Header.Get("Accept-Language"), config.DefaultLanguageType)
					config.SessionManager.Put(c.Request().Context(), asessions.ECHOSCS_OBJECTKEY_USER_SESSION, us)
				}
				// Set the "us" struct to this echo.Context so it is available elsewhere.
				c.Set(asessions.ECHOSCS_OBJECTKEY_USER_SESSION, &us)
			}

			c.Response().Before(func() {
				if config.SessionManager.Status(ctx) != scs.Unmodified {
					responseCookie := &http.Cookie{
						Name:     config.SessionManager.Cookie.Name,
						Path:     config.SessionManager.Cookie.Path,
						Domain:   config.SessionManager.Cookie.Domain,
						Secure:   config.SessionManager.Cookie.Secure,
						HttpOnly: config.SessionManager.Cookie.HttpOnly,
						SameSite: config.SessionManager.Cookie.SameSite,
					}

					switch config.SessionManager.Status(ctx) {
					case scs.Modified:
						token, _, err = config.SessionManager.Commit(ctx)
						if err != nil {
							panic(err)
						}

						responseCookie.Value = token

					case scs.Destroyed:
						responseCookie.Expires = time.Unix(1, 0)
						responseCookie.MaxAge = -1
					}

					c.SetCookie(responseCookie)
					addHeaderIfMissing(c.Response(), "Cache-Control", `no-cache="Set-Cookie"`)
					addHeaderIfMissing(c.Response(), "Vary", "Cookie")
				}
			})

			return next(c)
		}
	}
}

// addHeaderIfMissing adds a header to the response if it is not already set.
func addHeaderIfMissing(w http.ResponseWriter, key, value string) {
	if _, found := w.Header()[key]; !found {
		w.Header().Add(key, value)
	}
}
