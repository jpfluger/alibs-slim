package anetwork

import (
	"net"
	"testing"
)

// TestIsOutsidePortRange checks if the IsOutsidePortRange function correctly identifies invalid port numbers.
func TestIsOutsidePortRange(t *testing.T) {
	tests := []struct {
		port     int
		expected bool
	}{
		{0, true},
		{1, false},
		{80, false},
		{65534, false},
		{65535, true},
	}

	for _, test := range tests {
		if IsOutsidePortRange(test.port) != test.expected {
			t.Errorf("IsOutsidePortRange(%d) = %v, want %v", test.port, !test.expected, test.expected)
		}
	}
}

// TestGetExternalIPFromLocalInterfaces checks if the GetExternalIPFromLocalInterfaces function can return an IP address.
func TestGetExternalIPFromLocalInterfaces(t *testing.T) {
	ip, err := GetExternalIPFromLocalInterfaces()
	if err != nil {
		t.Errorf("GetExternalIPFromLocalInterfaces() error = %v", err)
	}
	if net.ParseIP(ip) == nil {
		t.Errorf("GetExternalIPFromLocalInterfaces() returned an invalid IP: %v", ip)
	}
}

// TestGetAllExternalIPFromLocalInterfaces checks if the GetAllExternalIPFromLocalInterfaces function can return a list of IP addresses.
func TestGetAllExternalIPFromLocalInterfaces(t *testing.T) {
	ips, err := GetAllExternalIPFromLocalInterfaces()
	if err != nil {
		t.Errorf("GetAllExternalIPFromLocalInterfaces() error = %v", err)
	}
	if len(ips) == 0 {
		t.Error("GetAllExternalIPFromLocalInterfaces() returned no IP addresses")
	}
}

// TestGetExternalIPByDNS checks if the GetExternalIPByDNS function can return an IP address using a DNS server.
func TestGetExternalIPByDNS(t *testing.T) {
	dns := "8.8.8.8" // Google's public DNS server
	ip, err := GetExternalIPByDNS(dns)
	if err != nil {
		t.Errorf("GetExternalIPByDNS(%s) error = %v", dns, err)
	}
	if net.ParseIP(ip.String()) == nil {
		t.Errorf("GetExternalIPByDNS(%s) returned an invalid IP: %v", dns, ip)
	}
}

// TestGetExternalIPByResolver checks if the GetExternalIPByResolver function can return an IP address using the system resolver or a fallback DNS server.
func TestGetExternalIPByResolver(t *testing.T) {
	dns := "8.8.8.8" // Fallback DNS server
	ip, err := GetExternalIPByResolver(dns)
	if err != nil {
		t.Errorf("GetExternalIPByResolver(%s) error = %v", dns, err)
	}
	if net.ParseIP(ip.String()) == nil {
		t.Errorf("GetExternalIPByResolver(%s) returned an invalid IP: %v", dns, ip)
	}
}
