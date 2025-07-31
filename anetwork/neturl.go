package anetwork

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ParseNetURL parses a string into a NetURL object.
func ParseNetURL(target string) (*NetURL, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("could not parse string to create NetURL: %v", err)
	}
	// Perform additional validation to ensure the URL has at least a scheme and host.
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("parsed URL is missing required components")
	}
	return &NetURL{URL: parsedURL}, nil
}

// MustParseNetURL attempts to parse a string into a NetURL object and returns an empty NetURL on error.
func MustParseNetURL(target string) *NetURL {
	parsedURL, err := ParseNetURL(target)
	if err != nil {
		return &NetURL{}
	}
	return parsedURL
}

// ParseNetURLNoError attempts to parse a string into a NetURL object and returns nil on error.
func ParseNetURLNoError(target string) *NetURL {
	parsedURL, err := ParseNetURL(target)
	if err != nil {
		return nil
	}
	return parsedURL
}

// GetUrlPathOrRoot returns the path of the URL or "/" if the path is empty.
func GetUrlPathOrRoot(u *url.URL) string {
	if u == nil || u.Path == "" {
		return "/"
	}
	return u.Path
}

// NetURL wraps the standard net/url URL type.
type NetURL struct {
	*url.URL
}

// GetSchemeHost returns the scheme and host of the URL as a string.
func (nu *NetURL) GetSchemeHost() string {
	return fmt.Sprintf("%s://%s", nu.Scheme, nu.Host)
}

// MarshalJSON implements the json.Marshaler interface for NetURL.
func (nu NetURL) MarshalJSON() ([]byte, error) {
	if nu.URL == nil {
		return []byte(`null`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, nu.URL.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for NetURL.
func (nu *NetURL) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 {
		nu.URL = nil
		return nil
	}

	su := strings.Trim(string(b), `"`)
	if su == "null" {
		return nil
	}

	parsedURL, err := url.Parse(su)
	if err != nil {
		return fmt.Errorf("NetURL.UnmarshalJSON could not parse URL: %v", err)
	}
	nu.URL = parsedURL
	return nil
}

// NewUrlJoinPath joins the given path elements to the existing URL and returns the new URL as a string.
func (nu *NetURL) NewUrlJoinPath(elem ...string) (string, error) {
	if nu.URL == nil {
		return "", fmt.Errorf("url is empty")
	}

	if strings.TrimSpace(nu.URL.Host) == "" {
		return "", fmt.Errorf("url is empty")
	}

	// Create a copy of the original URL to avoid modifying it
	newURL := *nu.URL

	if len(elem) == 0 {
		return newURL.String(), nil
	}

	// Join the new path elements
	joinedPath, err := url.JoinPath(newURL.Path, elem...)
	if err != nil {
		return "", fmt.Errorf("could not join path: %v", err)
	}
	newURL.Path = joinedPath

	return newURL.String(), nil
}

// MustJoin joins the given path elements to the existing URL and returns the new URL as a string.
func (nu *NetURL) MustJoin(elem ...string) string {
	newURL, err := nu.NewUrlJoinPath(elem...)
	if err != nil {
		return ""
	}
	return newURL
}

// SplitHostPortWithDefaults splits the host and port, applying default ports based on the scheme if necessary.
func (nu *NetURL) SplitHostPortWithDefaults(applyPortByScheme bool) (host string, port string, err error) {
	if nu.Port() == "" && applyPortByScheme {
		port = map[string]string{"https": "443", "http": "80"}[nu.Scheme]
		return nu.Host, port, nil
	}
	host, port, err = net.SplitHostPort(nu.Host)
	return
}

// IsHttps checks if the URL scheme is HTTPS.
func (nu *NetURL) IsHttps() bool {
	return nu.URL != nil && nu.URL.Scheme == "https"
}

// IsUrl checks if the NetURL is a valid URL.
func (nu *NetURL) IsUrl() bool {
	return nu != nil && nu.URL != nil && nu.URL.Scheme != "" && nu.URL.Host != ""
}

// GetPortForStartServer returns the port to start a server, including the domain if specified.
func (nu *NetURL) GetPortForStartServer(includeDomain bool) string {
	if includeDomain {
		return nu.Host
	}
	port := map[bool]string{true: ":443", false: ":80"}[nu.IsHttps()]
	if nu.Port() != "" {
		port = ":" + nu.Port()
	}
	return port
}

// GetPortNoDefaultHasColon returns the port with a colon prefix if set, otherwise an empty string.
func (nu *NetURL) GetPortNoDefaultHasColon() string {
	if nu.Port() != "" {
		return ":" + nu.Port()
	}
	return ""
}

// String returns the string representation of the NetURL.
func (nu *NetURL) String() string {
	if nu == nil || nu.URL == nil {
		return ""
	}
	return nu.URL.String()
}

// IsHostIPAsString checks if the host is an IP address and returns it as a string.
func (nu *NetURL) IsHostIPAsString() (string, bool) {
	if nu == nil || nu.URL == nil {
		return "", false
	}
	addr := net.ParseIP(nu.URL.Host)
	return nu.URL.Host, addr != nil
}

// IsReachable checks if the URL is reachable by attempting a TCP connection.
func (nu *NetURL) IsReachable() (bool, error) {
	if nu.URL == nil {
		return false, fmt.Errorf("the URL is nil")
	}
	host, port, err := nu.SplitHostPortWithDefaults(true)
	if err != nil {
		return false, fmt.Errorf("failed to parse host:port for URL '%s': %v", nu.String(), err)
	}
	if host == "" || port == "" {
		return false, fmt.Errorf("unknown URL '%s'", nu.String())
	}
	_, err = net.DialTimeout("tcp", net.JoinHostPort(host, port), 1*time.Second)
	return err == nil, err
}

// GetListenerKey returns the listener key in the format of "host:port".
func (nu *NetURL) GetListenerKey() string {
	if nu.URL == nil || !nu.IsUrl() {
		return ""
	}
	host, port := nu.Hostname(), nu.GetPortForStartServer(false)
	if strings.HasPrefix(port, ":") {
		port = strings.TrimPrefix(port, ":")
	}
	return net.JoinHostPort(host, port)
}

// IsUrlReachableWithPing checks if the URL is reachable by pinging it a specified number of times.
func (nu *NetURL) IsUrlReachableWithPing(maxPing int, sleepDuration int) bool {
	if nu == nil {
		return false
	}
	for counter := 0; maxPing == 0 || counter < maxPing; counter++ {
		time.Sleep(time.Duration(sleepDuration) * time.Second)
		if ok, _ := nu.IsReachable(); ok {
			return true
		}
	}
	return false
}

// IsUrlUnreachableWithPing checks if the URL is unreachable by pinging it a specified number of times.
func (nu *NetURL) IsUrlUnreachableWithPing(maxPing int, sleepDuration int) bool {
	if nu == nil {
		return false
	}
	for counter := 0; maxPing == 0 || counter < maxPing; counter++ {
		time.Sleep(time.Duration(sleepDuration) * time.Second)
		if ok, _ := nu.IsReachable(); !ok {
			return true
		}
	}
	return false
}

// ReplaceWithPort replaces the current port of the URL with the specified port.
func (nu *NetURL) ReplaceWithPort(port int) *NetURL {
	if nu == nil || !nu.IsUrl() {
		return nil
	}
	domain := nu.Hostname()
	if port == 0 {
		port = map[bool]int{true: 443, false: 80}[nu.IsHttps()]
	}
	newURL := fmt.Sprintf("%s://%s:%d", nu.Scheme, domain, port)
	return ParseNetURLNoError(newURL)
}

// Copy creates a copy of the NetURL.
func (nu *NetURL) Copy() *NetURL {
	if nu == nil {
		return nil
	}
	newURL, _ := url.Parse(nu.URL.String())
	return &NetURL{URL: newURL}
}

// CleanUrl removes the port from the URL if it is a default port (80 or 443).
func (nu *NetURL) CleanUrl() error {
	if nu == nil || !nu.IsUrl() {
		return fmt.Errorf("URL is invalid")
	}
	portInt, err := nu.GetPortInt()
	if err != nil {
		return err
	}
	if portInt == 443 || portInt == 80 {
		nu.URL = &url.URL{Scheme: nu.Scheme, Host: nu.Hostname()}
	}
	return nil
}

// GetPortInt extracts the port number as an integer from the URL.
func (nu *NetURL) GetPortInt() (int, error) {
	if nu == nil || !nu.IsUrl() {
		return 0, fmt.Errorf("url is invalid")
	}
	portStr := nu.Port()
	if portStr == "" {
		return map[bool]int{true: 443, false: 80}[nu.IsHttps()], nil
	}
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert port to int (port=%s): %v", portStr, err)
	}
	if IsOutsidePortRange(portInt) {
		return 0, fmt.Errorf("port number out of range")
	}
	return portInt, nil
}

type NetURLs []NetURL
