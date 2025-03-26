package ahttp

import (
	"encoding/base64"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/alog"
	"github.com/jpfluger/alibs-slim/azb"
	"net/http"
	"time"
)

// Constants for cookie names used in flash messages.
const (
	COOKIE_FLASH_FORGOT_LOGIN   = azb.ZBType("fm-forgot-login")
	COOKIE_FLASH_RESET_PASSWORD = azb.ZBType("fm-reset-password")
)

// FlashMessageData holds the data to be displayed in a flash message to the user.
type FlashMessageData struct {
	Path     string `json:"path,omitempty" form:"path"`         // The path where the flash message originated.
	Action   string `json:"action,omitempty" form:"action"`     // The action that triggered the flash message.
	Username string `json:"username,omitempty" form:"username"` // The username involved in the action.
}

// FlashMessage represents the structure of a flash message.
type FlashMessage struct {
	Title   string `json:"title" form:"title"`     // The title of the flash message.
	Lead    string `json:"lead" form:"lead"`       // The lead text of the flash message.
	Message string `json:"message" form:"message"` // The main message content.
}

// RedirectFlashMessage sets a flash message cookie and redirects the user.
func RedirectFlashMessage(c echo.Context, cookieName azb.ZBType, fmData *FlashMessageData, targetUrl string) error {
	// Marshal the FlashMessageData into JSON.
	bFMData, err := json.Marshal(fmData)
	if err != nil {
		// Log the error using the application's logger.
		alog.LOGGER(alog.LOGGER_APP).Err(err).Str("fn", "redirect-flash-message").Str("cookie", cookieName.String()).Msg("failed fmData json marshal")
		return err
	}

	// Create a new cookie with the marshaled data.
	flash := &http.Cookie{
		Name:     cookieName.String(),
		Path:     "/",
		Value:    base64.URLEncoding.EncodeToString(bFMData),
		Expires:  time.Now().Add(30 * time.Second), // Set the cookie to expire in 30 seconds.
		HttpOnly: true,
	}

	// Set the cookie in the response.
	c.SetCookie(flash)

	// Redirect the user to the target URL.
	return RedirectCheckJSON(c, targetUrl)
}
