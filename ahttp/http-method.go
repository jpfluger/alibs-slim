package ahttp

import (
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// HttpMethod represents the HTTP methods used in requests.
type HttpMethod string

// Constants for the various HTTP methods.
const (
	HTTPMETHOD_GET    HttpMethod = "GET"
	HTTPMETHOD_POST   HttpMethod = "POST"
	HTTPMETHOD_PUT    HttpMethod = "PUT"
	HTTPMETHOD_DELETE HttpMethod = "DELETE"
	HTTPMETHOD_PATCH  HttpMethod = "PATCH"
)

// IsEmpty checks if the HttpMethod is empty after trimming whitespace.
func (hm HttpMethod) IsEmpty() bool {
	return strings.TrimSpace(string(hm)) == ""
}

// TrimSpace returns a new HttpMethod with leading and trailing whitespace removed.
func (hm HttpMethod) TrimSpace() HttpMethod {
	return HttpMethod(strings.TrimSpace(string(hm)))
}

// String converts the HttpMethod to a string.
func (hm HttpMethod) String() string {
	return string(hm)
}

// ToStringTrimUpper trims whitespace from the HttpMethod and converts it to uppercase.
func (hm HttpMethod) ToStringTrimUpper() string {
	return autils.ToStringTrimUpper(hm.String())
}

// TrimSpaceToUpper trims whitespace from the HttpMethod and converts it to uppercase.
func (hm HttpMethod) TrimSpaceToUpper() HttpMethod {
	return HttpMethod(autils.ToStringTrimUpper(hm.String()))
}

// HttpMethods is a slice of HttpMethod.
type HttpMethods []HttpMethod

// Find searches for a target HttpMethod in the slice and returns it if found.
func (hms HttpMethods) Find(target HttpMethod) HttpMethod {
	if hms == nil || len(hms) == 0 {
		return ""
	}
	for _, hm := range hms {
		if hm == target {
			return hm
		}
	}
	return ""
}

// Has checks if a target HttpMethod is present in the slice.
func (hms HttpMethods) Has(target HttpMethod) bool {
	return hms.Find(target) != ""
}
