package ahttp

import (
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/arob"
	"net/http"
	"strings"
)

// CheckMIMEType defines an enumeration for different MIME types to check against.
type CheckMIMEType int

const (
	// CHECK_MIME_TYPE_NONE indicates no MIME type is expected.
	CHECK_MIME_TYPE_NONE CheckMIMEType = iota
	// CHECK_MIME_TYPE_JSON indicates a JSON MIME type.
	CHECK_MIME_TYPE_JSON
	// CHECK_MIME_TYPE_XML indicates an XML MIME type.
	CHECK_MIME_TYPE_XML
	// CHECK_MIME_TYPE_HTML indicates an HTML MIME type.
	CHECK_MIME_TYPE_HTML
)

// MIMETYPE_JSON_NOUTF is used to specify the JSON MIME type without character set encoding.
const MIMETYPE_JSON_NOUTF = `application/json`

// GetRequestContentType retrieves the Content-Type header from the request.
func GetRequestContentType(c echo.Context) string {
	return c.Request().Header.Get(echo.HeaderContentType)
}

// IsRequestContentType checks if the request's Content-Type matches the specified CheckMIMEType.
func IsRequestContentType(c echo.Context, checkMimeType CheckMIMEType) bool {
	ctype := GetRequestContentType(c)
	if ctype == "" {
		return checkMimeType == CHECK_MIME_TYPE_NONE
	}

	switch {
	case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
		return checkMimeType == CHECK_MIME_TYPE_JSON
	case strings.HasPrefix(ctype, echo.MIMETextHTML), strings.HasPrefix(ctype, echo.MIMETextHTMLCharsetUTF8):
		return checkMimeType == CHECK_MIME_TYPE_HTML
	case strings.HasPrefix(ctype, echo.MIMEApplicationXML), strings.HasPrefix(ctype, echo.MIMETextXML):
		return checkMimeType == CHECK_MIME_TYPE_XML
	default:
		return false
	}
}

// RedirectCheckJSON redirects the client to the specified URL or returns a JSON response with the URL if the request is JSON.
func RedirectCheckJSON(c echo.Context, url string) error {
	if IsRequestContentType(c, CHECK_MIME_TYPE_JSON) {
		// Respond with JSON if the request expects a JSON response.
		return c.JSON(http.StatusOK, arob.NewROBWithRedirect(url))
	}
	// Otherwise, perform a standard HTTP redirect.
	return c.Redirect(http.StatusFound, url)
}

// DetectMimeTypeSendMessage sends a message with the appropriate MIME type.
// If the code is not provided, it defaults to http.StatusOK.
func DetectMimeTypeSendMessage(c echo.Context, code int, message string, isSysError bool) error {
	if code == 0 {
		code = http.StatusOK
	}
	return DetectMimeTypeSendMessageWithCode(c, code, message, isSysError)
}

// DetectMimeTypeSendMessageWithCode sends a message with the appropriate MIME type and status code.
func DetectMimeTypeSendMessageWithCode(c echo.Context, code int, message string, isSysError bool) error {
	if IsRequestContentType(c, CHECK_MIME_TYPE_JSON) {
		// If the request content type is JSON, respond with JSON.
		return c.JSON(code, arob.NewROBMessageWithOptionSysError(arob.ROBMessage(message), isSysError))
	} else if IsRequestContentType(c, CHECK_MIME_TYPE_XML) {
		// If the request content type is XML, respond with XML.
		return c.XML(code, arob.NewROBMessageWithOptionSysError(arob.ROBMessage(message), isSysError))
	}
	// Default response is a plain string.
	return c.String(code, message)
}

// DetectMimeTypeSendError sends an error message with the appropriate MIME type.
// If the code is not provided, it defaults to http.StatusInternalServerError.
func DetectMimeTypeSendError(c echo.Context, code int, err error) error {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	return DetectMimeTypeSendErrorWithCode(c, code, err)
}

// DetectMimeTypeSendErrorWithCode sends an error message with the appropriate MIME type and status code.
func DetectMimeTypeSendErrorWithCode(c echo.Context, code int, err error) error {
	if IsRequestContentType(c, CHECK_MIME_TYPE_JSON) {
		// If the request content type is JSON, respond with JSON.
		return c.JSON(code, arob.NewROBMessageWithOptionSysError(arob.ROBMessage(err.Error()), true))
	} else if IsRequestContentType(c, CHECK_MIME_TYPE_XML) {
		// If the request content type is XML, respond with XML.
		return c.XML(code, arob.NewROBMessageWithOptionSysError(arob.ROBMessage(err.Error()), true))
	}
	// Default response is a plain string.
	return c.String(code, err.Error())
}

// APIJSONSendROBError sends a JSON response with a status of "OK" (http.StatusOK) and an ROB error type.
// It uses APIJSONSendROBStatusWithOptions to return a response with an error message.
// Parameters:
// - c: The echo.Context object, which carries the request and response context.
// - err: The error object that will be included in the response as the error detail.
func APIJSONSendROBError(c echo.Context, err error) error {
	return APIJSONSendROBStatusWithOptions(c, http.StatusOK, arob.ROBTYPE_ERROR, "error", nil, err)
}

// APIJSONSendROBErrorWithCode sends a JSON response with a custom HTTP status code and an ROB error type.
// It allows specifying the HTTP status code in the response, useful for sending different error codes.
// Parameters:
// - c: The echo.Context object.
// - code: The custom HTTP status code to be returned (e.g., 400 for Bad Request).
// - err: The error object that will be included in the response as the error detail.
func APIJSONSendROBErrorWithCode(c echo.Context, code int, err error) error {
	return APIJSONSendROBStatusWithOptions(c, code, arob.ROBTYPE_ERROR, "error", nil, err)
}

// APIJSONSendROB sends a JSON response with data wrapped in an ROB record.
// It is used to send a successful response (http.StatusOK) with additional data.
// Parameters:
// - c: The echo.Context object.
// - data: The data to be included in the response.
func APIJSONSendROB(c echo.Context, data interface{}) error {
	return APIJSONSendROBStatusWithOptions(c, http.StatusOK, "", "", data, nil)
}

// APIJSONSendROBWithCode sends a JSON response with a custom HTTP status code and data wrapped in an ROB record.
// This allows specifying a custom status code for the response.
// Parameters:
// - c: The echo.Context object.
// - code: The custom HTTP status code.
// - data: The data to be included in the response.
func APIJSONSendROBWithCode(c echo.Context, code int, data interface{}) error {
	return APIJSONSendROBStatusWithOptions(c, code, "", "", data, nil)
}

// APIJSONSendROBInfoOk sends a JSON response with an "OK" status and ROB info type.
// It is a convenience function for sending a simple "OK" message in an ROB response.
// Parameters:
// - c: The echo.Context object.
func APIJSONSendROBInfoOk(c echo.Context) error {
	return APIJSONSendROBStatusWithOptions(c, http.StatusOK, arob.ROBTYPE_INFO, "", nil, nil)
}

// APIJSONSendROBInfoStatus sends a JSON response with a status of http.StatusOK (200) and a custom status message.
// It includes a custom status message in the ROB response.
// Parameters:
// - c: The echo.Context object.
// - status: The custom status message to be included in the response.
func APIJSONSendROBInfoStatus(c echo.Context, status string) error {
	return APIJSONSendROBStatusWithOptions(c, http.StatusOK, arob.ROBTYPE_INFO, status, nil, nil)
}

// APIJSONSendROBInfoRecs sends a JSON response with a status of http.StatusOK (200) and records wrapped in an ROB record.
// It is used to send a response with both a custom status message and additional data records.
// Parameters:
// - c: The echo.Context object.
// - status: The custom status message.
// - recs: The records to be included in the response.
func APIJSONSendROBInfoRecs(c echo.Context, status string, recs interface{}) error {
	return APIJSONSendROBStatusWithOptions(c, http.StatusOK, arob.ROBTYPE_INFO, status, recs, nil)
}

// APIJSONSendROBStatusWithOptions is a utility function that sends a JSON response with flexible options.
// It allows sending responses with different status codes, ROB types, status messages, data, and error details.
// Parameters:
// - c: The echo.Context object.
// - code: The HTTP status code (e.g., 200 for OK, 400 for Bad Request).
// - statusType: The type of ROB message (e.g., ROBTYPE_INFO, ROBTYPE_ERROR).
// - status: The custom status message (if empty, it defaults to the standard status text for the HTTP code).
// - data: Optional data to include in the response.
// - err: Optional error object to include in the response (used if non-nil).
func APIJSONSendROBStatusWithOptions(c echo.Context, code int, statusType arob.ROBType, status string, data interface{}, err error) error {
	// If no HTTP status code is provided, default to http.StatusOK (200).
	if code == 0 {
		code = http.StatusOK
	}
	// If no ROB type is provided, default to ROBTYPE_INFO if there's no error, otherwise ROBTYPE_ERROR.
	if statusType.IsEmpty() {
		if err == nil {
			statusType = arob.ROBTYPE_INFO
		} else {
			statusType = arob.ROBTYPE_ERROR
		}
	}
	// If no custom status message is provided, use the default HTTP status text for the given code.
	if status == "" {
		status = http.StatusText(code)
	}
	// If an error is provided, send a response with the error details.
	if err != nil {
		return c.JSON(code, arob.NewROBWithStatusError(statusType, status, err))
	}
	// If no data is provided, send a response with just the status message.
	if data == nil {
		return c.JSON(code, arob.NewROBWithStatus(statusType, status))
	}
	// Otherwise, send a response with both the status message and the data records.
	return c.JSON(code, arob.NewROBWithStatusRecs(statusType, status, data))
}
