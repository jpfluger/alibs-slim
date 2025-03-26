package ahttp

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// IHTTPErrorHandler defines the interface for handling HTTP errors.
type IHTTPErrorHandler interface {
	GetErr() error
	GetContext() echo.Context
	GetLogger() echo.Logger
	GetIsOnDebug() bool
	GetHttpCode() int
	GetHttpMessage() string
	SetHttpCode(httpCode int)
	SetHttpMessage(httpMessage string)
	HandleResponse() error // Implementation required at a higher level.
}

// HTTPErrorHandlerBase provides a basic implementation of IHTTPErrorHandler.
type HTTPErrorHandlerBase struct {
	err       error
	c         echo.Context
	logger    echo.Logger
	isOnDebug bool

	httpCode    int
	httpMessage string
}

// NewHTTPErrorHandlerBase creates a new instance of HTTPErrorHandlerBase.
func NewHTTPErrorHandlerBase(err error, c echo.Context, logger echo.Logger, isOnDebug bool) *HTTPErrorHandlerBase {
	return &HTTPErrorHandlerBase{
		err:       err,
		c:         c,
		logger:    logger,
		isOnDebug: isOnDebug,
	}
}

// Getters and setters for HTTPErrorHandlerBase.
func (he *HTTPErrorHandlerBase) GetErr() error                     { return he.err }
func (he *HTTPErrorHandlerBase) GetContext() echo.Context          { return he.c }
func (he *HTTPErrorHandlerBase) GetLogger() echo.Logger            { return he.logger }
func (he *HTTPErrorHandlerBase) GetIsOnDebug() bool                { return he.isOnDebug }
func (he *HTTPErrorHandlerBase) GetHttpCode() int                  { return he.httpCode }
func (he *HTTPErrorHandlerBase) GetHttpMessage() string            { return he.httpMessage }
func (he *HTTPErrorHandlerBase) SetHttpCode(httpCode int)          { he.httpCode = httpCode }
func (he *HTTPErrorHandlerBase) SetHttpMessage(httpMessage string) { he.httpMessage = httpMessage }

// DefaultHTTPErrorHandler handles HTTP errors by sending a JSON response with the status code.
// It is a variation of echo.DefaultHTTPErrorHandler.
func DefaultHTTPErrorHandler(options IHTTPErrorHandler) {
	// If the response has already been committed, return early.
	if options.GetContext().Response().Committed {
		return
	}

	err := options.GetErr()
	he, ok := err.(*echo.HTTPError)
	if ok {
		// If there is an internal error, use it.
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		// Create a new HTTP error with the internal server error status.
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	// Set the HTTP code and message based on the error.
	options.SetHttpCode(he.Code)
	if m, ok := he.Message.(string); ok {
		if options.GetIsOnDebug() {
			options.SetHttpMessage(err.Error())
		} else {
			options.SetHttpMessage(m)
		}
	}

	// Send the response.
	if options.GetContext().Request().Method == http.MethodHead {
		err = options.GetContext().NoContent(he.Code)
	} else {
		err = options.HandleResponse()
	}
	if err != nil {
		options.GetLogger().Error(err)
	}
}
