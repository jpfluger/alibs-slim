package anetwork

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ParseNetURL parses and validates a URL with open-ended, authority-based schemes,
// and normalizes file: URLs to their canonical forms.
// Accepted:
//   - file:///abs/path          (empty authority; absolute path required)
//   - file://host/abs/path      (authority present; absolute path required)
//   - Any other scheme WHEN written in //authority form (e.g., smb://host, rdp://host, arl://host)
//
// Rejected:
//   - Opaque/no-authority forms (e.g., "mailto:", "ssh:host")
//   - file:// (no path), file:/ (no path), file:/rel/path (not absolute)
func ParseNetURL(target string) (*NetURL, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("could not parse string to create NetURL: %w", err)
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("parsed URL is missing scheme")
	}

	switch strings.ToLower(u.Scheme) {
	case "file":
		u, err = normalizeFileURL(u, target)
		if err != nil {
			return nil, err
		}

	default:
		rest := strings.TrimPrefix(target, u.Scheme+":")
		hasAuthority := strings.HasPrefix(rest, "//")
		if !hasAuthority {
			return nil, fmt.Errorf("%s URL must use //authority form (e.g., %s://host/...)", u.Scheme, u.Scheme)
		}
		if u.Host == "" {
			return nil, fmt.Errorf("%s URL requires a host", u.Scheme)
		}
	}

	return &NetURL{URL: u}, nil
}

func normalizeFileURL(u *url.URL, original string) (*url.URL, error) {
	rest := strings.TrimPrefix(original, u.Scheme+":")
	hasAuthority := strings.HasPrefix(rest, "//")

	if !hasAuthority {
		// Allow file:/abs but force file:///abs
		if u.Path == "" || !strings.HasPrefix(u.Path, "/") || u.Path == "/" {
			return nil, fmt.Errorf("file URL must be absolute and not just root (e.g. file:///path)")
		}
		return &url.URL{Scheme: "file", Path: u.Path}, nil
	}

	// Authority form
	if u.Host == "" {
		if u.Path == "" || !strings.HasPrefix(u.Path, "/") || u.Path == "/" {
			return nil, fmt.Errorf("file URL without host must be absolute and not just root (file:///path)")
		}
		return &url.URL{Scheme: "file", Path: u.Path}, nil
	}

	// file://host/abs/path
	if u.Path == "" || !strings.HasPrefix(u.Path, "/") || u.Path == "/" {
		return nil, fmt.Errorf("file URL with host must have absolute non-root path like file://host/path")
	}
	return &url.URL{Scheme: "file", Host: u.Host, Path: u.Path}, nil
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

// String safely returns the string representation of the URL.
func (nu *NetURL) String() string {
	if nu == nil || nu.URL == nil {
		return ""
	}
	return nu.URL.String()
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

// IsEmpty checks if the NetURL is nil or represents an empty or whitespace-only URL string.
func (nu *NetURL) IsEmpty() bool {
	if nu.URL == nil {
		return true
	}
	return strings.TrimSpace(nu.URL.String()) == ""
}

// IsFile checks if the NetURL has a "file" scheme (case-insensitive) and a non-empty path.
// Returns false if the URL is nil or the scheme is not "file".
func (nu *NetURL) IsFile() bool {
	if nu.URL == nil {
		return false
	}
	return strings.EqualFold(nu.URL.Scheme, "file") && nu.URL.Path != "" && nu.URL.Path != "/"
}

// IsHttps checks if the NetURL has an "https" scheme (case-insensitive) and a non-empty host.
// Returns false if the URL is nil or the scheme is not "https".
func (nu *NetURL) IsHttps() bool {
	if nu.URL == nil {
		return false
	}
	return strings.EqualFold(nu.URL.Scheme, "https") && nu.URL.Host != ""
}

// IsHttpProtocol checks if the NetURL has a "http" or "https" scheme (case-insensitive) and a non-empty host.
func (nu *NetURL) IsHttpProtocol() bool {
	if nu.URL == nil {
		return false
	}
	scheme := strings.ToLower(nu.URL.Scheme)
	if scheme != "http" && scheme != "https" {
		return false
	}
	return nu.URL.Host != ""
}

// IsURL reports whether the NetURL is a supported, usable URL.
// - file: requires non-empty, non-root absolute path
// - everything else: requires non-empty host
func (nu *NetURL) IsUrl() bool {
	if nu == nil || nu.URL == nil {
		return false
	}
	scheme := strings.ToLower(nu.URL.Scheme)
	if scheme == "" {
		return false
	}
	if scheme == "file" {
		// Accept only absolute, non-root path for file URLs
		return nu.URL.Path != "" && nu.URL.Path != "/"
	}
	// Authority-style network schemes: host required
	return nu.URL.Host != ""
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

// IsHostIP checks if the host is an IP address.
func (nu *NetURL) IsHostIP() bool {
	if nu == nil || nu.URL == nil {
		return false
	}
	addr := net.ParseIP(nu.URL.Host)
	return addr != nil
}

// IsHostIPAsString checks if the host is an IP address as a string.
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
	newURL, _ := url.Parse(nu.String())
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

// FindNextOpenPort tries to bind sequentially starting at `port` until it
// finds an available TCP port (supports dual-stack IPv4/IPv6).
//
// * port <= 0    → 49152 (first IANA-registered ephemeral port)
// * returns the open port number or an error if none found.
func FindNextOpenPort(port int) (int, error) {
	if port < NETPORT_MIN || port > NETPORT_MAX {
		port = NETPORT_EPHEMERAL
	}

	for p := port; p <= NETPORT_MAX; p++ {
		addr := fmt.Sprintf(":%d", p)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			// Port is in use – keep searching.
			continue
		}
		_ = ln.Close() // Immediately free it.
		return p, nil
	}

	return 0, fmt.Errorf("no open port found in range %d-%d", port, NETPORT_MAX)
}

// WithNextOpenPort clones the current NetURL, discovers the next open port
// (starting at startPort) on the supplied host, and returns the updated copy
// together with the chosen port.
//
// When startPort <= 0 the search starts at the first ephemeral port (49152).
func (nu *NetURL) WithNextOpenPort(startPort int) (*NetURL, int, error) {
	if nu == nil || !nu.IsUrl() {
		return nil, 0, fmt.Errorf("nil or invalid NetURL")
	}

	openPort, err := FindNextOpenPort(startPort)
	if err != nil {
		return nil, 0, err
	}

	// Leverage the existing helper to build a fresh NetURL with the new port.
	return nu.ReplaceWithPort(openPort), openPort, nil
}

// isAcceptableWebHost accept only IPs, localhost, or hosts containing a dot.
func isAcceptableWebHost(h string) bool {
	if h == "" {
		return false
	}
	if h == "localhost" {
		return true
	}
	if net.ParseIP(h) != nil {
		return true
	}
	return strings.Contains(h, ".")
}

// NormalizeURL adds https:// to schemeless web inputs and defers validation/normalization to ParseNetURL.
// Non-web schemes (file:, smb:, rdp:, arl:, etc.) are returned unchanged.
func NormalizeURL(u NetURL) NetURL {
	if u.IsEmpty() {
		return NetURL{}
	}
	orig := strings.TrimSpace(u.String())
	if orig == "" {
		return NetURL{}
	}

	scheme := strings.ToLower(u.Scheme)

	switch scheme {
	case "http", "https":
		parsed, err := ParseNetURL(orig)
		if err != nil {
			return NetURL{}
		}
		// STRICT web host check (fixes the failing tests)
		if !isAcceptableWebHost(parsed.Hostname()) {
			return NetURL{}
		}
		return *parsed

	case "":
		// Build https://<domain><path> from schemeless input
		domain := strings.TrimPrefix(strings.TrimSpace(u.Host), "//")
		path := u.Path
		rawQuery := u.RawQuery
		fragment := u.Fragment

		if domain == "" {
			full := strings.TrimSpace(u.Path)
			if full == "" {
				return NetURL{}
			}
			parts := strings.SplitN(full, "/", 2)
			domain = parts[0]
			if len(parts) > 1 {
				path = "/" + parts[1]
			} else {
				path = ""
			}
		}
		if domain == "" {
			return NetURL{}
		}

		b := &url.URL{
			Scheme:   "https",
			Host:     domain,
			Path:     path,
			RawQuery: rawQuery,
			Fragment: fragment,
		}
		parsed, err := ParseNetURL(b.String())
		if err != nil {
			return NetURL{}
		}
		// STRICT web host check (fixes the failing tests)
		if !isAcceptableWebHost(parsed.Hostname()) {
			return NetURL{}
		}
		return *parsed

	default:
		// Non-web scheme: leave unchanged (ParseNetURL handles its own validation/normalization)
		return u
	}
}

type NetURLs []NetURL

func (nus NetURLs) Clean(enforceHTTPs bool) NetURLs {
	arr := NetURLs{}
	for _, nu := range nus {
		if nu.IsEmpty() {
			continue
		}
		if enforceHTTPs && !nu.IsHttps() {
			continue
		}
		arr = append(arr, nu)
	}
	return arr
}
