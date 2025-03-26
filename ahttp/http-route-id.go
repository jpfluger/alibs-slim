package ahttp

import "strings"

// HttpRouteId is a custom type for route identifiers.
type HttpRouteId string

// IsEmpty checks if the HttpRouteId is empty after trimming whitespace.
func (rt HttpRouteId) IsEmpty() bool {
	return strings.TrimSpace(string(rt)) == ""
}

// TrimSpace returns a new HttpRouteId with leading and trailing whitespace removed.
func (rt HttpRouteId) TrimSpace() HttpRouteId {
	return HttpRouteId(strings.TrimSpace(string(rt)))
}

// String converts the HttpRouteId to a string.
func (rt HttpRouteId) String() string {
	return string(rt)
}

// Define constants for various route identifiers using the NAME TYPE = VALUE format.
const (
	RPAGE_ROOT                     HttpRouteId = "RPAGE_ROOT"
	RPAGE_ROOT_HOME                HttpRouteId = "RPAGE_ROOT_HOME"
	RPAGE_ROOT_NOT_FOUND           HttpRouteId = "RPAGE_ROOT_NOT_FOUND"
	RPAGE_ROOT_SERVICE_UNAVAILABLE HttpRouteId = "RPAGE_ROOT_SERVICE_UNAVAILABLE"
	RPAGE_ROOT_FORBIDDEN           HttpRouteId = "RPAGE_ROOT_FORBIDDEN"
	RPAGE_US_ACTION                HttpRouteId = "RPAGE_US_ACTION"
	RPAGE_US_ACTION_POST           HttpRouteId = "RPAGE_US_ACTION_POST"
	RPAGE_US_ACTION_UNAVAILABLE    HttpRouteId = "RPAGE_US_ACTION_UNAVAILABLE"
	RPAGE_ROOT_MAINTENANCE         HttpRouteId = "RPAGE_ROOT_MAINTENANCE"
	ROUTE_REDIRECT_DEFAULT         HttpRouteId = "ROUTE_REDIRECT_DEFAULT"
	ROUTE_REDIRECT_AUTH_ERROR      HttpRouteId = "ROUTE_REDIRECT_AUTH_ERROR"
	RPAGE_LOGIN                    HttpRouteId = "RPAGE_LOGIN"
	RPAGE_LOGIN_SUBMIT             HttpRouteId = "RPAGE_LOGIN_SUBMIT"
	RPAGE_LOGOUT                   HttpRouteId = "RPAGE_LOGOUT"
	RPAGE_FORGOT_LOGIN             HttpRouteId = "RPAGE_FORGOT_LOGIN"
	RPAGE_FORGOT_LOGIN_SUBMIT      HttpRouteId = "RPAGE_FORGOT_LOGIN_SUBMIT"
	RPAGE_FLASH                    HttpRouteId = "RPAGE_FLASH"
	RPAGE_JWT_LINK                 HttpRouteId = "RPAGE_JWT_LINK"
	RPAGE_LEGAL_SUMMARY            HttpRouteId = "RPAGE_LEGAL_SUMMARY"
	RPAGE_LEGAL_TERMS              HttpRouteId = "RPAGE_LEGAL_TERMS"
	RPAGE_LEGAL_PRIVACY            HttpRouteId = "RPAGE_LEGAL_PRIVACY"
	RPAGE_LEGAL_COOKIES            HttpRouteId = "RPAGE_LEGAL_COOKIES"
	RPAGE_LEGAL_DMCA               HttpRouteId = "RPAGE_LEGAL_DMCA"
	RPAGE_SYSTEM_SGLOBALS          HttpRouteId = "RPAGE_SYSTEM_SGLOBALS"
	RPAGE_SYSTEM_SGLOBALS_MOD      HttpRouteId = "RPAGE_SYSTEM_SGLOBALS_MOD"
	RPAGE_SYSTEM_SDOCS             HttpRouteId = "RPAGE_SYSTEM_SDOCS"
	RPAGE_SYSTEM_SDOCS_MOD         HttpRouteId = "RPAGE_SYSTEM_SDOCS_MOD"
	RPAGE_SYSTEM_SMAILERS          HttpRouteId = "RPAGE_SYSTEM_SMAILERS"
	RPAGE_SYSTEM_SMAILERS_MOD      HttpRouteId = "RPAGE_SYSTEM_SMAILERS_MOD"
	RPAGE_CONTACT_CDESK            HttpRouteId = "RPAGE_CONTACT_CDESK"
	RPAGE_CONTACT_CDESK_MOD        HttpRouteId = "RPAGE_CONTACT_CDESK_MOD"
	RPAGE_USER_DASHBOARD           HttpRouteId = "RPAGE_USER_DASHBOARD"
	RPAGE_USER_DASHBOARD_MOD       HttpRouteId = "RPAGE_USER_DASHBOARD_MOD"
	RPAGE_USER_PROFILE             HttpRouteId = "RPAGE_USER_PROFILE"
	RPAGE_USER_PROFILE_MOD         HttpRouteId = "RPAGE_USER_PROFILE_MOD"
	ROUTE_REDIRECT_CUSTOM          HttpRouteId = "ROUTE_REDIRECT_CUSTOM"
	RAPI_PING                      HttpRouteId = "RAPI_PING"
	RAPI_UNITTEST_SHUTDOWNS        HttpRouteId = "RAPI_UNITTEST_SHUTDOWNS"
)
