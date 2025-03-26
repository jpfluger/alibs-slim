package asessions

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUserSession tests the creation of UserSessionPerm and its status.
func TestUserSession(t *testing.T) {
	// Define a function that always returns true for testing purposes.
	fn := func(target ILoginSessionPerm) bool {
		return true
	}

	// Create a new UserSessionPerm with default status and test its status.
	us := NewUserSessionPerm()
	assert.Equal(t, LOGIN_SESSION_STATUS_NONE, us.Status, "Default status should be LOGIN_SESSION_STATUS_NONE")
	assert.True(t, fn(us), "Function should return true for the new user session")

	// Create a new UserSessionPerm with a specific status and test its status.
	us = NewUserSessionPermWithLoginStatus(LOGIN_SESSION_STATUS_OK)
	assert.Equal(t, LOGIN_SESSION_STATUS_OK, us.Status, "Status should be LOGIN_SESSION_STATUS_OK after setting it")
	assert.True(t, fn(us), "Function should return true for the user session with set status")
}

// TestCastLoginSessionPermFromEchoContext tests the retrieval of ILoginSessionPerm from Echo's context.
func TestCastLoginSessionPermFromEchoContext(t *testing.T) {
	e := echo.New()
	req := e.AcquireContext().Request()
	rec := e.AcquireContext().Response().Writer
	c := e.NewContext(req, rec)

	// Simulate storing a UserSessionPerm in Echo's context.
	expectedSession := NewUserSessionPerm()
	c.Set(ECHOSCS_OBJECTKEY_USER_SESSION, expectedSession)

	// Cast the session from the context and assert it is correct.
	session := CastLoginSessionPermFromEchoContext(c)
	assert.Equal(t, expectedSession, session, "Retrieved session should match the stored session.")
}

// TestCastUserSessionPermFromEchoContext tests the retrieval of *UserSessionPerm from Echo's context.
func TestCastUserSessionPermFromEchoContext(t *testing.T) {
	e := echo.New()
	req := e.AcquireContext().Request()
	rec := e.AcquireContext().Response().Writer
	c := e.NewContext(req, rec)

	// Simulate storing a UserSessionPerm in Echo's context.
	expectedSession := NewUserSessionPerm()
	c.Set(ECHOSCS_OBJECTKEY_USER_SESSION, expectedSession)

	// Cast the session from the context and assert it is correct.
	session := CastUserSessionPermFromEchoContext(c)
	assert.Equal(t, expectedSession, session, "Retrieved session should match the stored session.")
}

// TestCastUserSessionPermFromILoginSession tests the casting of ILoginSessionPerm to *UserSessionPerm.
func TestCastUserSessionPermFromILoginSession(t *testing.T) {
	expectedSession := NewUserSessionPerm()
	session := CastUserSessionPermFromILoginSession(expectedSession)
	assert.Equal(t, expectedSession, session, "Casted session should match the original session.")
}

func TestUserSessionPermMethods(t *testing.T) {
	session := NewUserSessionPermWithLoginStatus(LOGIN_SESSION_STATUS_OK)

	// Test HasPerm method.
	session.Perms = NewPermSetByBits("read", PERM_R)

	// Test HasPerm for no permissions
	assert.False(t, session.HasPerm(Perm{key: "read", value: &PermValue{}}), "HasPerm should return false for no 'read' permission.")
	assert.False(t, session.HasPerm(Perm{key: "read", value: nil}), "HasPerm should return false for nil 'read' permission.")

	// Test HasPerm for valid permissions
	assert.True(t, session.HasPerm(Perm{key: "read", value: &PermValue{value: PERM_R}}), "HasPerm should return true for 'read:R' permission.")

	// Test HasPermS for string-based permissions
	assert.True(t, session.HasPermS("read:R"), "HasPermS should return true for 'read:R' permission.")

	// Test HasPermB for keyBits-based permissions
	assert.True(t, session.HasPermB("read:2"), "HasPermB should return true for 'read:R' permission.")

	// Test HasPermBV for specific bit permissions
	assert.True(t, session.HasPermBV("read", PERM_R), "HasPermBV should return true for 'read' with bit PERM_R.")

	// Test HasPermSV for key and string value-based permissions
	assert.True(t, session.HasPermSV("read", "R"), "HasPermSV should return true for 'read:R' permission.")

	// Test invalid HasPermSV cases
	assert.False(t, session.HasPermSV("write", "R"), "HasPermSV should return false for non-existent 'write:R' permission.")
	assert.False(t, session.HasPermSV("read", "X"), "HasPermSV should return false for invalid 'read:X' permission.")
}
