package ahttp

import (
	"fmt"
	"mime"
	"net/http"
	"net/http/httputil"
	"strings"
)

// ISON_UNITTESTS_WAIT_USER_SHUTDOWN specifies the wait time in seconds for user-initiated service shutdown during unit tests.
var ISON_UNITTESTS_WAIT_USER_SHUTDOWN = 0

// ISON_UNITTESTS_UPDOWN_SECRET is the secret key required to authorize a unit test shutdown.
var ISON_UNITTESTS_UPDOWN_SECRET = ""

// GetIsOnUnitTests checks if unit tests are configured to wait for user shutdown.
func GetIsOnUnitTests() bool {
	return ISON_UNITTESTS_WAIT_USER_SHUTDOWN != 0
}

// GetIsOnUnitTestsHasSecret checks if a secret key is set for unit test shutdown.
func GetIsOnUnitTestsHasSecret() bool {
	return ISON_UNITTESTS_UPDOWN_SECRET != ""
}

// HasContentType checks if the request has the specified MIME type in its Content-Type header.
func HasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}

// DumpHeaderRequest prints the HTTP request headers and body for debugging purposes.
func DumpHeaderRequest(req *http.Request, dumpBody bool) {
	requestDump, err := httputil.DumpRequest(req, dumpBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(requestDump))
}

// DumpHeaderResponse prints the HTTP response headers and body for debugging purposes.
func DumpHeaderResponse(res *http.Response, dumpBody bool) {
	resDump, err := httputil.DumpResponse(res, dumpBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resDump))
}

// ExtractHttpRouteId extracts the HttpRouteId from a string if it is prefixed with "hrt:".
func ExtractHttpRouteId(name string) (HttpRouteId, bool) {
	if name == "" || !strings.HasPrefix(name, "hrt:") {
		return "", false
	}
	parts := strings.Split(name, ":")
	if len(parts) != 2 {
		return "", false
	}
	routeId := HttpRouteId(strings.TrimSpace(parts[1]))
	if routeId.IsEmpty() {
		return "", false
	}
	return routeId, true
}

// JoinUrl concatenates the root URL with a path, ensuring proper formatting.
func JoinUrl(urlRoot string, urlPath string) string {
	urlRoot = strings.TrimSuffix(urlRoot, "/")
	urlPath = strings.TrimPrefix(urlPath, "/")
	return fmt.Sprintf("%s/%s", urlRoot, urlPath)
}
