package amidware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/asessions"
)

// mockAuthenticateProvisioner implements the IAuthenticateProvisioner interface for testing.
type mockAuthenticateProvisioner struct{}

func (m *mockAuthenticateProvisioner) GetUrlNoLogin() string {
	return "/no-login"
}

func (m *mockAuthenticateProvisioner) GetUrlInvalidPerms() string {
	return "/invalid-perms"
}

func (m *mockAuthenticateProvisioner) LogAuthError(c echo.Context, err error) {
	// Log the error or handle it as needed for your tests.
}

// TestAuthenticateConfig tests the AuthenticateConfig middleware.
func TestAuthenticateConfig_ValidPerm(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock session and set it in the Echo context.
	session := asessions.NewUserSessionPerm()
	session.Status = asessions.LOGIN_SESSION_STATUS_OK
	session.Perms = asessions.NewPermSetByPair("read", "XLCRUD")

	// Add the session to the context.
	c.Set(asessions.ECHOSCS_OBJECTKEY_USER_SESSION, session)

	// Define permissions required for the test.
	requiredPerms := asessions.NewPermSetByBits("read", asessions.PERM_R)

	// Initialize the middleware with a mock provisioner.
	middleware := NewAuthenticatePermConfig(requiredPerms, &mockAuthenticateProvisioner{})

	// Define a handler to represent the next handler in the middleware chain.
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Test the middleware with a valid session.
	h := middleware(handler)
	err := h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}

// TestAuthenticateConfig tests the AuthenticateConfig middleware.
func TestAuthenticateConfig_InvalidPerm(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock session and set it in the Echo context.
	session := asessions.NewUserSessionPerm()
	session.Status = asessions.LOGIN_SESSION_STATUS_OK
	session.Perms = asessions.NewPermSetByPair("read", "X")

	// Add the session to the context.
	c.Set(asessions.ECHOSCS_OBJECTKEY_USER_SESSION, session)

	// Define permissions required for the test.
	requiredPerms := asessions.NewPermSetByBits("read", asessions.PERM_R)

	// Initialize the middleware with a mock provisioner.
	middleware := NewAuthenticatePermConfig(requiredPerms, &mockAuthenticateProvisioner{})

	// Define a handler to represent the next handler in the middleware chain.
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Test the middleware with a valid session.
	h := middleware(handler)
	err := h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "", rec.Body.String())

	// Add more test cases to cover different scenarios, such as:
	// - No session present in the context.
	// - Session with status LOGIN_SESSION_STATUS_NONE.
	// - Skipping the middleware.
	// - Redirects when no login or invalid permissions are detected.
}
