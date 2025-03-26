package amidware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/asessions"
)

// mockProvisioner implements the IActionRoutesProvisioner interface for testing purposes.
// mockProvisioner implements the IActionRoutesProvisioner interface for testing purposes.
type mockProvisioner struct{}

func (m *mockProvisioner) MatchesWhitelistActionPath(targetPath string) bool {
	// Mock logic for testing
	return targetPath == "/whitelisted-path"
}

func (m *mockProvisioner) LogError(c echo.Context, err error) {
	// Mock logging, you can simulate logging or leave it empty for testing purposes
	c.Logger().Error(err)
}

func (m *mockProvisioner) GetUrlNotFoundActions() string {
	// Mock return value for "not found" actions URL
	return "/notfound"
}

func (m *mockProvisioner) IsMaintenanceMode() bool {
	// Mock maintenance mode status
	return false
}

func (m *mockProvisioner) GetUrlMaintenance() string {
	// Mock return value for maintenance mode URL
	return "/maintenance"
}

// TestActionRoutes tests the ActionRoutes middleware.
func TestActionRoutes(t *testing.T) {

	session := asessions.NewUserSessionPerm()
	session.Status = asessions.LOGIN_SESSION_STATUS_OK
	//actionKey := asessions.ActionKey("2factor")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Add the session to the context.
	c.Set(asessions.ECHOSCS_OBJECTKEY_USER_SESSION, session)

	// Define your ActionKeyUrls and other necessary setup for the test.
	actionKeyUrls := testActionKeyUrls
	//actionKeyUrls := asessions.ActionKeyUrls{
	//	// Populate with test data.
	//}

	// Initialize the middleware with a mock provisioner.
	middleware := NewActionRoutes(actionKeyUrls, &mockProvisioner{})

	// Define a handler that represents the next handler in the middleware chain.
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Test the middleware.
	h := middleware(handler)
	err := h(c)

	// Assert that the middleware behaves as expected.
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())

	// Add more test cases as needed to cover different scenarios.
}

// testActionKeyUrls is an array of ActionKeyUrl suitable for testing.
var testActionKeyUrls = asessions.ActionKeyUrls{
	{
		Key:         "home",
		Url:         "/home",
		CheckPrefix: false,
	},
	{
		Key:         "profile",
		Url:         "/user/profile",
		CheckPrefix: true,
	},
	{
		Key:         "settings",
		Url:         "/user/settings",
		CheckPrefix: false,
	},
	{
		Key:         "logout",
		Url:         "/logout",
		CheckPrefix: false,
	},
	{
		Key:         "dashboard",
		Url:         "/dashboard",
		CheckPrefix: true,
	},
	// Add more ActionKeyUrl instances as needed for comprehensive testing.
}

// validateActionRoutesConfig checks the provided configuration for any missing or invalid values.
//func validateActionRoutesConfig(config *ActionRoutesConfig) error {
//	// Check if the configuration is nil
//	if config == nil {
//		return errors.New("config cannot be nil")
//	}
//
//	// Check if the provisioner is nil
//	if config.Provisioner == nil {
//		return errors.New("Provisioner cannot be nil")
//	}
//
//	// Check if ActionKeyUrls is nil or empty
//	if config.ActionKeyUrls == nil || len(config.ActionKeyUrls) == 0 {
//		return errors.New("ActionKeyUrls cannot be empty")
//	}
//
//	// Validate each ActionKeyUrl in the list
//	for _, aku := range config.ActionKeyUrls {
//		if aku == nil {
//			return errors.New("ActionKeyUrl entry cannot be nil")
//		}
//		if aku.Key.IsEmpty() {
//			return errors.New("ActionKeyUrl.Key cannot be empty")
//		}
//		if strings.TrimSpace(aku.Url) == "" {
//			return errors.New("ActionKeyUrl.Url cannot be empty")
//		}
//	}
//
//	return nil
//}

// TestValidateActionRoutesConfig tests the validation of ActionRoutesConfig.
func TestValidateActionRoutesConfig(t *testing.T) {
	// Test with a nil config
	err := validateActionRoutesConfig(nil)
	assert.Error(t, err)
	assert.Equal(t, "config is nil", err.Error())

	// Test with nil Provisioner
	err = validateActionRoutesConfig(&ActionRoutesConfig{
		ActionKeyUrls: []*asessions.ActionKeyUrl{},
		Provisioner:   nil,
	})
	assert.Error(t, err)
	assert.Equal(t, "Provisioner is nil", err.Error())

	// Test with nil ActionKeyUrls
	err = validateActionRoutesConfig(&ActionRoutesConfig{
		ActionKeyUrls: nil,
		Provisioner:   &mockProvisioner{},
	})
	assert.Error(t, err)
	assert.Equal(t, "ActionKeyUrls is empty", err.Error())

	// Test with an empty ActionKeyUrls
	err = validateActionRoutesConfig(&ActionRoutesConfig{
		ActionKeyUrls: []*asessions.ActionKeyUrl{},
		Provisioner:   &mockProvisioner{},
	})
	assert.Error(t, err)
	assert.Equal(t, "ActionKeyUrls is empty", err.Error())

	// Test with an invalid ActionKeyUrl (empty Key and Url)
	err = validateActionRoutesConfig(&ActionRoutesConfig{
		ActionKeyUrls: []*asessions.ActionKeyUrl{
			{
				Key:         asessions.ActionKey(""), // Assuming ActionKey{}.IsEmpty() is true
				Url:         "",
				CheckPrefix: false,
			},
		},
		Provisioner: &mockProvisioner{},
	})
	assert.Error(t, err)
	assert.Equal(t, "ActionKeyUrls.Key is empty", err.Error())

	// Test with a valid config
	err = validateActionRoutesConfig(&ActionRoutesConfig{
		ActionKeyUrls: []*asessions.ActionKeyUrl{
			{
				Key:         asessions.ActionKey("test-key"), // Assuming ActionKey{}.IsEmpty() is false
				Url:         "/test-url",
				CheckPrefix: false,
			},
		},
		Provisioner: &mockProvisioner{},
	})
	assert.NoError(t, err)
}
