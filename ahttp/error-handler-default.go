package ahttp

import (
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/arob"
	"github.com/jpfluger/alibs-slim/asessions"
	"net/http"
)

// HTTPErrorHandlerDefault is a default error handler that requires a PSC
// and uses fixed pages, like "status.gohtml". This struct may be too inflexible to your needs.
// If this is the case, adapt it in your higher-level application.
type HTTPErrorHandlerDefault struct {
	HTTPErrorHandlerBase
	newPD NewPageData
}

func NewHTTPErrorHandlerDefault(err error, c echo.Context, logger echo.Logger, isOnDebug bool, newPD NewPageData) *HTTPErrorHandlerDefault {
	return &HTTPErrorHandlerDefault{
		HTTPErrorHandlerBase: *NewHTTPErrorHandlerBase(err, c, logger, isOnDebug),
		newPD:                newPD,
	}
}

func (he *HTTPErrorHandlerDefault) HandleResponse() error {
	if IsRequestContentType(he.GetContext(), CHECK_MIME_TYPE_JSON) {
		message := he.GetHttpMessage()
		if message == "" {
			message = http.StatusText(he.GetHttpCode())
		}
		return he.GetContext().JSON(he.GetHttpCode(), arob.NewROBWithError(arob.ROBERRORFIELD_SYSTEM, arob.ROBMessage(message)))
	}

	routeId := RPAGE_ROOT_SERVICE_UNAVAILABLE
	switch he.GetHttpCode() {
	case http.StatusNotFound:
		routeId = RPAGE_ROOT_NOT_FOUND
	case http.StatusForbidden:
		routeId = RPAGE_ROOT_FORBIDDEN
	default:
		break
	}

	c := he.GetContext()
	return c.Render(http.StatusOK, "status.gohtml", he.newPD(PSC().MustUrl(routeId), "Unavailable", asessions.CastLoginSessionPermFromEchoContext(c), &PageStatusDefault{HTTPCode: he.GetHttpCode()}))
}

// RHStatus returns http.StatusOK. If JSON is detected, it automatically creates and returns a rob error object.
func RHStatus(c echo.Context, code int, routeId HttpRouteId, newPD NewPageData) error {

	// routeId can be used to return different pages.
	if routeId.IsEmpty() {
		routeId = RPAGE_ROOT_SERVICE_UNAVAILABLE
	}

	pageStatusDefault := &PageStatusDefault{HTTPCode: code, RouteId: routeId}

	if IsRequestContentType(c, CHECK_MIME_TYPE_JSON) {
		message := pageStatusDefault.StatusMessage()
		return c.JSON(http.StatusOK, arob.NewROBWithError(arob.ROBERRORFIELD_SYSTEM, arob.ROBMessage(message)))
	}

	return c.Render(http.StatusOK, "status.gohtml", newPD(PSC().MustUrl(routeId), "Status", asessions.CastLoginSessionPermFromEchoContext(c), pageStatusDefault))
}

// PageStatusDefault is a default struct to display status for web pages.
// This may be too limiting to your situation, especially for multi-language support.
type PageStatusDefault struct {
	HTTPCode     int
	MessageTitle string
	Message      string
	RouteId      HttpRouteId
}

func (p *PageStatusDefault) StatusTitle() string {
	if p.MessageTitle != "" {
		return p.MessageTitle
	} else if p.HTTPCode == http.StatusNotFound {
		return "Resource not found"
	} else if p.HTTPCode == http.StatusInternalServerError {
		return "Service unavailable"
	} else if p.HTTPCode == http.StatusForbidden {
		return "Forbidden"
	} else if p.RouteId == RPAGE_ROOT_MAINTENANCE {
		return "Maintenance"
	} else {
		return "Service Unavailable"
	}
}

func (p *PageStatusDefault) StatusMessage() string {
	if p.Message != "" {
		return p.Message
	} else if p.HTTPCode == http.StatusNotFound {
		return "The requested resource could not be found. Try the link below to get back on track."
	} else if p.HTTPCode == http.StatusInternalServerError {
		return "An unexpected condition was encountered. Try the link below to a known serviceable page."
	} else if p.HTTPCode == http.StatusForbidden {
		return "You do not have permission to access this resource."
	} else if p.RouteId == RPAGE_ROOT_MAINTENANCE {
		return "The site is currently under maintenance. Please check back later."
	} else {
		return "An unexpected condition was encountered. Try the link below to a known serviceable page."
	}
}
