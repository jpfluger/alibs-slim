package ahttp

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/azb"
	"io"
)

// RHBind binds the request body to the provided data interface for POST requests.
func RHBind(c echo.Context, din interface{}) error {
	// Only process POST requests.
	if c.Request().Method != "POST" {
		return nil
	}
	// Bind the request body to the provided data interface.
	if err := c.Bind(din); err != nil {
		return err
	}
	return nil
}

// RHBindNopCloser binds the request body to the provided data interface for POST requests
// and restores the request body to allow for subsequent reads.
func RHBindNopCloser(c echo.Context, din interface{}) error {
	// Only process POST requests.
	if c.Request().Method != "POST" {
		return nil
	}

	// Read the request body and store it for potential reuse.
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request().Body)
	}

	// Restore the io.ReadCloser to its original state.
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Bind the request body to the provided data interface.
	if err := c.Bind(din); err != nil {
		return err
	}

	// Restore the io.ReadCloser again to its original state.
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

// RHBindListPostOnly binds the request body to the provided data interface only for POST requests.
func RHBindListPostOnly(c echo.Context, din azb.IDINPaginate) error {
	// Only process POST requests.
	if c.Request().Method != "POST" {
		return nil
	}
	return RHBindList(c, din)
}

// RHBindList binds the request body to the provided data interface and validates it.
func RHBindList(c echo.Context, din azb.IDINPaginate) error {
	// Check if the data interface is nil.
	if din == nil {
		return fmt.Errorf("din is nil")
	}

	// Bind the request body to the provided data interface.
	if err := c.Bind(&din); err != nil {
		return err
	}

	// Validate the bound data.
	if err := din.Validate(); err != nil {
		return fmt.Errorf("validate failed; %v", err)
	}

	return nil
}

// RHBindListQuery binds query parameters to the provided data interface and validates it.
func RHBindListQuery(c echo.Context, din azb.IDINPaginate) error {
	// Check if the data interface is nil.
	if din == nil {
		return fmt.Errorf("din is nil")
	}

	// Bind the query parameters to the provided data interface.
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, din); err != nil {
		return err
	}

	// Validate the bound data.
	if err := din.Validate(); err != nil {
		return fmt.Errorf("validate failed; %v", err)
	}

	return nil
}
