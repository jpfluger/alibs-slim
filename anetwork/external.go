package anetwork

import (
	"fmt"
	"net"
)

const (
	NETPORT_MIN       = 1     // Minimum valid port number, 0 is reserved
	NETPORT_MAX       = 65534 // Maximum valid port number, 65535 is reserved
	NETPORT_EPHEMERAL = 49152
)

// IsOutsidePortRange checks if a given port number is outside the valid range of 1 to 65535.
func IsOutsidePortRange(portInt int) bool {
	return portInt < NETPORT_MIN || portInt > NETPORT_MAX
}

// GetExternalIPFromLocalInterfaces attempts to find the external IP address from the local network interfaces.
func GetExternalIPFromLocalInterfaces() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // Skip down or loopback interfaces
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip != nil && ip.To4() != nil {
				return ip.String(), nil // Return the first non-loopback IPv4 address
			}
		}
	}
	return "", fmt.Errorf("cannot determine the external IP from local interfaces")
}

// InterfaceIP holds the IP addresses associated with a network interface.
type InterfaceIP struct {
	InterfaceName string   `json:"interfaceName"`
	IPv4          []net.IP `json:"ipv4"`
	IPv6          []net.IP `json:"ipv6"`
}

// InterfaceIPs is a slice of pointers to InterfaceIP.
type InterfaceIPs []*InterfaceIP

// GetAllExternalIPFromLocalInterfaces retrieves all external IP addresses from local network interfaces.
func GetAllExternalIPFromLocalInterfaces() (InterfaceIPs, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var myIntIPs InterfaceIPs
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // Skip down or loopback interfaces
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue // Ignore errors and try the next interface
		}

		currentIntIP := &InterfaceIP{InterfaceName: iface.Name}
		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}

			if ipv4 := ip.To4(); ipv4 != nil {
				currentIntIP.IPv4 = append(currentIntIP.IPv4, ipv4)
			} else if ipv6 := ip.To16(); ipv6 != nil {
				currentIntIP.IPv6 = append(currentIntIP.IPv6, ipv6)
			}
		}

		if len(currentIntIP.IPv4) > 0 || len(currentIntIP.IPv6) > 0 {
			myIntIPs = append(myIntIPs, currentIntIP) // Add only if there are IP addresses
		}
	}
	return myIntIPs, nil
}

// GetExternalIPByDNS attempts to find the external IP address by connecting to a DNS server.
func GetExternalIPByDNS(dns string) (net.IP, error) {
	if dns == "" {
		return nil, fmt.Errorf("dns parameter is empty")
	}

	conn, err := net.Dial("udp", dns+":80")
	if err != nil {
		return nil, fmt.Errorf("unable to find default outbound ip: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// GetExternalIPByResolver attempts to find the external IP address using the system resolver or a fallback DNS server.
func GetExternalIPByResolver(dnsIfResolverEmpty string) (net.IP, error) {
	resolv, err := GetResolv()
	if err == nil {
		for _, ns := range resolv.Nameservers {
			if ip, err := GetExternalIPByDNS(ns); err == nil {
				return ip, nil
			}
		}
	}

	if ip, err := GetExternalIPByDNS(dnsIfResolverEmpty); err == nil {
		return ip, nil
	}

	return nil, fmt.Errorf("cannot determine ip address")
}

// extractIP extracts the IP address from the given Addr.
func extractIP(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	}
	return nil
}
