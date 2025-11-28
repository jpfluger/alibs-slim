package aclient_http

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockXML struct {
	Message string `xml:"message"`
}

func setupEchoServer() *echo.Echo {
	e := echo.New()

	e.HEAD("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, JSON!"})
	})

	e.GET("/xml", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXML)
		return c.XML(http.StatusOK, MockXML{Message: "Hello, XML!"})
	})

	e.POST("/json", func(c echo.Context) error {
		var payload map[string]string
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
		}
		return c.JSON(http.StatusOK, payload)
	})

	e.POST("/xml", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXML)

		// Read the request body as a string
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.XML(http.StatusBadRequest, MockXML{Message: "Invalid payload"})
		}
		bodyString := string(bodyBytes)
		fmt.Println("Request Body:", bodyString)

		// Decode the request body into the payload
		var payload MockXML
		if err := xml.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&payload); err != nil && err != io.EOF {
			return c.XML(http.StatusBadRequest, MockXML{Message: "Invalid payload"})
		}
		return c.XML(http.StatusOK, payload)
	})

	return e
}

func TestAClientHTTP_GetJSON(t *testing.T) {
	e := setupEchoServer()
	server := httptest.NewServer(e)
	defer server.Close()

	client, err := NewAClientHTTP(server.URL)
	assert.NoError(t, err)

	hob := NewHOBGet("/json")
	var result map[string]string

	err = client.GetJSON(hob, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, JSON!", result["message"])
}

func TestAClientHTTP_GetXML(t *testing.T) {
	e := setupEchoServer()
	server := httptest.NewServer(e)
	defer server.Close()

	client, err := NewAClientHTTP(server.URL)
	assert.NoError(t, err)

	hob := NewHOBGet("/xml")
	var result MockXML

	err = client.GetXML(hob, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, XML!", result.Message)
}

func TestAClientHTTP_PostJSON(t *testing.T) {
	e := setupEchoServer()
	server := httptest.NewServer(e)
	defer server.Close()

	client, err := NewAClientHTTP(server.URL)
	assert.NoError(t, err)

	hob := NewHOBWithJSON("/json", map[string]string{"message": "Hello, JSON!"})
	var result map[string]string

	err = client.PostJSON(hob, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, JSON!", result["message"])
}

func TestAClientHTTP_PostXML(t *testing.T) {
	e := setupEchoServer()
	server := httptest.NewServer(e)
	defer server.Close()

	client, err := NewAClientHTTP(server.URL)
	assert.NoError(t, err)

	hob := NewHOBWithXML("/xml", MockXML{Message: "Hello, XML!"})
	var result MockXML

	err = client.PostXML(hob, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, XML!", result.Message)
}
