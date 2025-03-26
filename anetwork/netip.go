package anetwork

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// NetIP holds individual IPs and subnets.
type NetIP struct {
	IP     net.IP
	Subnet *net.IPNet
}

// NetIPs is a slice of NetIP.
type NetIPs []NetIP

// Contains checks if the ipTarget is in the list of IPs or subnets.
func (ns NetIPs) Contains(ipTarget string) bool {
	targetIP := net.ParseIP(ipTarget)
	if targetIP == nil {
		return false
	}

	for _, netIP := range ns {
		if netIP.IP != nil && netIP.IP.Equal(targetIP) {
			return true
		}
		if netIP.Subnet != nil && netIP.Subnet.Contains(targetIP) {
			return true
		}
	}
	return false
}

// ToNetIP parses a slice of IP addresses and CIDRs, returning a slice of NetIP and an error.
func ToNetIP(netIPs []string) (NetIPs, error) {
	var parsedIPs NetIPs

	for _, ipStr := range netIPs {
		ipStr = strings.TrimSpace(ipStr)
		if strings.HasSuffix(ipStr, "/32") {
			ipStr = strings.TrimSuffix(ipStr, "/32")
		}

		if _, ipNet, err := net.ParseCIDR(ipStr); err == nil {
			// If it's a valid CIDR, add the IPNet
			parsedIPs = append(parsedIPs, NetIP{Subnet: ipNet})
		} else if ip := net.ParseIP(ipStr); ip != nil {
			// If it's a valid IP, add it
			parsedIPs = append(parsedIPs, NetIP{IP: ip})
		} else {
			return nil, fmt.Errorf("invalid IP address or CIDR: %s", ipStr)
		}
	}

	return parsedIPs, nil
}

// ConvertIPToFilename converts an IP address to a safe filename format with optional extension.
func ConvertIPToFilename(ip string, ext string) string {
	// Replace non-alphanumeric characters (like `:`, `.`) with `_`
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	safeIP := re.ReplaceAllString(ip, "_")

	ext = strings.TrimSpace(ext)
	if strings.HasPrefix(ext, ".") {
		ext = strings.TrimPrefix(ext, ".")
	}
	if strings.TrimSpace(ext) == "" {
		return safeIP
	}
	// Append extension
	return fmt.Sprintf("%s.%s", safeIP, ext)
}

// ConvertIPLogsToFilenames converts each IP in IPLogs to a safe filename with optional extension.
func ConvertIPLogsToFilenames(ips []string, ext string) []string {
	filenames := make([]string, len(ips))
	for i, ip := range ips {
		filenames[i] = ConvertIPToFilename(ip, ext)
	}
	return filenames
}
