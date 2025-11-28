package amidware

import (
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jpfluger/alibs-slim/asessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Register the session types with gob for serialization.
	gob.Register(asessions.UserSessionPerm{})
}

func TestSessionUserSession(t *testing.T) {

	gob.Register(asessions.UserSessionPerm{})

	var sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	e := echo.New()

	// Call /put to set the message in the session manager
	req := httptest.NewRequest(http.MethodGet, "/put", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	session := SCSLoadAndSave(sessionManager, true)

	am := &asessions.PermSet{}
	am.SetPerm(asessions.NewPerm("access:X"))
	am.SetPerm(asessions.NewPerm("admin:XCRUD"))
	am.SetPerm(asessions.NewPerm("employee:"))

	h := session(func(c echo.Context) error {
		sessionManager.Put(c.Request().Context(), "message", "Hello from a session!")

		us := asessions.CastUserSessionPermFromEchoContext(c)
		us.DisplayName = "UserSession works too!"
		us.Perms = *am
		sessionManager.Put(c.Request().Context(), asessions.ECHOSCS_OBJECTKEY_USER_SESSION, us)

		return c.String(http.StatusOK, "")
	})

	h(c)

	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, 1, len(rec.Result().Cookies()))

	sessionCookie := rec.Result().Cookies()[0]

	assert.Equal(t, "session", sessionCookie.Name)

	// Make a request to /get to see if the message is still there
	req = httptest.NewRequest(http.MethodGet, "/get", nil)
	req.Header.Set("Cookie", sessionCookie.Name+"="+sessionCookie.Value)

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	h = session(func(c echo.Context) error {
		msg := sessionManager.GetString(c.Request().Context(), "message")

		us := asessions.CastLoginSessionPermFromEchoContext(c)
		// Note: us.GetRUIDString() always displays "" because it RUID is a pointer.
		//       We are testing to ensure RUID is in fact empty!
		msg += " " + us.GetDisplayName()

		if us.GetRecordUserIdentity().String() == "" {
			msg += " RUID is empty!"
		}

		isError := false
		fnCheckError := func(expectedTruth bool, expectedKV string) {
			if isError {
				return // no point in evaluating something broken
			}
			if expectedTruth {
				if !am.HasPermS(expectedKV) {
					isError = true
					t.Logf("MatchesOne failed expectedTruth '%t' where perm=%s", expectedTruth, expectedKV)
				}
			} else {
				if am.HasPermS(expectedKV) {
					isError = true
					t.Logf("MatchesOne failed expectedTruth '%t' where perm=%s", expectedTruth, expectedKV)
				}
			}
		}

		fnCheckError(false, "nothing-test:X")
		fnCheckError(true, "access:X")
		fnCheckError(true, "admin:C")
		fnCheckError(false, "employee:C")

		if !isError {
			msg += " And Perms works too!"
		}

		return c.String(http.StatusOK, msg)
	})

	h(c)

	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, "Hello from a session! UserSession works too! RUID is empty! And Perms works too!", rec.Body.String())
}

func TestSession(t *testing.T) {
	// Register the session types with gob for serialization.
	gob.Register(asessions.UserSessionPerm{})

	// Initialize the session manager.
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	e := echo.New()

	// Create a middleware handler with the session manager.
	sessionMiddleware := SCSLoadAndSave(sessionManager, false)

	var cName, cValue string

	// Test setting a value in the session.
	t.Run("SetSessionValue", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/put", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := sessionMiddleware(func(c echo.Context) error {
			sessionManager.Put(c.Request().Context(), "message", "Hello from a session!")
			return c.String(http.StatusOK, "")
		})

		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Len(t, rec.Result().Cookies(), 1)

		sessionCookie := rec.Result().Cookies()[0]
		assert.Equal(t, sessionManager.Cookie.Name, sessionCookie.Name)
		cName = sessionCookie.Name
		cValue = sessionCookie.Value
	})

	// Test retrieving a value from the session.
	t.Run("GetSessionValue", func(t *testing.T) {
		// Create a new request for the /get route and include the session cookie.
		req := httptest.NewRequest(http.MethodGet, "/get", nil)
		req.Header.Set("Cookie", cName+"="+cValue)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := sessionMiddleware(func(c echo.Context) error {
			// Retrieve the message from the session.
			msg := sessionManager.GetString(c.Request().Context(), "message")
			return c.String(http.StatusOK, msg)
		})

		// Execute the handler and perform assertions.
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello from a session!", rec.Body.String())
	})
}

//
//func TestSessionUserSession2(t *testing.T) {
//
//	// Initialize the session manager.
//	sessionManager := scs.New()
//	sessionManager.Lifetime = 24 * time.Hour
//
//	e := echo.New()
//
//	// Create a middleware handler with the session manager.
//	sessionMiddleware := SCSLoadAndSave(sessionManager, true)
//
//	// Test setting and retrieving a custom user session.
//	t.Run("UserSessionManagement", func(t *testing.T) {
//		req := httptest.NewRequest(http.MethodGet, "/put", nil)
//		rec := httptest.NewRecorder()
//		c := e.NewContext(req, rec)
//
//		perms := asessions.PermSet{}
//		perms.SetByString("access:X")
//		perms.SetByString("admin:XCRUD")
//		perms.SetByString("employee:")
//
//		handler := sessionMiddleware(func(c echo.Context) error {
//			// Set the message in the session.
//			sessionManager.Put(c.Request().Context(), "message", "Hello from a session!")
//			// Create and set the custom user session.
//			us := asessions.NewUserSessionPerm()
//			us.DisplayName = "UserSession works too!"
//			us.Perms = perms
//			sessionManager.Put(c.Request().Context(), asessions.ECHOSCS_OBJECTKEY_USER_SESSION, *us)
//			return c.String(http.StatusOK, "")
//		})
//
//		assert.NoError(t, handler(c))
//		assert.Equal(t, http.StatusOK, rec.Code)
//		assert.Len(t, rec.Result().Cookies(), 1)
//
//		sessionCookie := rec.Result().Cookies()[0]
//		assert.Equal(t, sessionManager.Cookie.Name, sessionCookie.Name)
//
//		req = httptest.NewRequest(http.MethodGet, "/get", nil)
//		req.Header.Set("Cookie", sessionCookie.Name+"="+sessionCookie.Value)
//		rec = httptest.NewRecorder()
//		c = e.NewContext(req, rec)
//
//		handler = sessionMiddleware(func(c echo.Context) error {
//			// Retrieve the message from the session.
//			msg := sessionManager.GetString(c.Request().Context(), "message")
//			// Retrieve the custom user session.
//			us := asessions.CastUserSessionPermFromEchoContext(c)
//			msg += " " + us.DisplayName
//			return c.String(http.StatusOK, msg)
//		})
//
//		assert.NoError(t, handler(c))
//		assert.Equal(t, http.StatusOK, rec.Code)
//		assert.Equal(t, "Hello from a session! UserSession works too!", rec.Body.String())
//	})
//}
