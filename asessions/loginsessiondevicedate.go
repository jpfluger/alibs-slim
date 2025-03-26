package asessions

import (
	"fmt"
	"net"
	"time"
)

// LoginSessionDeviceDate represents the details of a device used in a login session.
type LoginSessionDeviceDate struct {
	Device string    `json:"device,omitempty"` // The device name or identifier.
	IP     string    `json:"ip,omitempty"`     // The IP address of the device as a string.
	Date   time.Time `json:"date,omitempty"`   // The date and time when the device was used.
}

// ParseIP attempts to parse the IP field and return a net.IP address.
// Returns an error if the IP field is empty or the format is invalid.
func (lsdd *LoginSessionDeviceDate) ParseIP() (net.IP, error) {
	if lsdd.IP == "" {
		return nil, fmt.Errorf("parse IP error: IP address is empty")
	}
	ip := net.ParseIP(lsdd.IP)
	if ip == nil {
		return nil, fmt.Errorf("parse IP error: invalid IP address format")
	}
	return ip, nil
}

// IsIPInSubnets checks if the IP of the LoginSessionDeviceDate is within any of the provided subnets.
// The subnets parameter is a slice of CIDR-formatted IP addresses or networks.
// Returns true if the IP is within any of the subnets, false otherwise.
func (lsdd *LoginSessionDeviceDate) IsIPInSubnets(subnets []string) (bool, error) {
	deviceIP, err := lsdd.ParseIP()
	if err != nil {
		return false, err
	}

	for _, subnet := range subnets {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			return false, fmt.Errorf("invalid subnet format: %v", err)
		}
		if ipNet.Contains(deviceIP) {
			return true, nil
		}
	}
	return false, nil
}

// IsIPInSubnetList checks if the IP of the LoginSessionDeviceDate is within any of the provided subnet list.
// The subnets parameter is a slice of *net.IPNet which are already parsed subnets.
// Returns true if the IP is within any of the subnets, false otherwise.
func (lsdd *LoginSessionDeviceDate) IsIPInSubnetList(subnets []*net.IPNet) (bool, error) {
	deviceIP, err := lsdd.ParseIP()
	if err != nil {
		return false, err
	}

	for _, ipNet := range subnets {
		if ipNet.Contains(deviceIP) {
			return true, nil
		}
	}
	return false, nil
}

// LoginSessionDeviceDates is a slice of pointers to LoginSessionDeviceDate.
// It holds the history of devices used in login sessions.
type LoginSessionDeviceDates []*LoginSessionDeviceDate
