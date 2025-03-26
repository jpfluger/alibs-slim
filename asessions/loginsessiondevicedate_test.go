package asessions

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestParseIP(t *testing.T) {
	// Valid IP address
	deviceDate := LoginSessionDeviceDate{IP: "192.168.1.1"}
	parsedIP, err := deviceDate.ParseIP()
	assert.NoError(t, err)
	assert.NotNil(t, parsedIP)

	// Invalid IP address
	deviceDate = LoginSessionDeviceDate{IP: "invalid_ip"}
	_, err = deviceDate.ParseIP()
	assert.Error(t, err)

	// Empty IP address
	deviceDate = LoginSessionDeviceDate{IP: ""}
	_, err = deviceDate.ParseIP()
	assert.Error(t, err)
}

func TestIsIPInSubnets(t *testing.T) {
	deviceDate := LoginSessionDeviceDate{IP: "192.168.1.1"}

	// IP is within the subnet
	subnets := []string{"192.168.1.0/24"}
	inSubnet, err := deviceDate.IsIPInSubnets(subnets)
	assert.NoError(t, err)
	assert.True(t, inSubnet)

	// IP is not within the subnet
	subnets = []string{"10.0.0.0/8"}
	inSubnet, err = deviceDate.IsIPInSubnets(subnets)
	assert.NoError(t, err)
	assert.False(t, inSubnet)

	// Invalid subnet format
	subnets = []string{"invalid_subnet"}
	_, err = deviceDate.IsIPInSubnets(subnets)
	assert.Error(t, err)
}

func TestIsIPInSubnetList(t *testing.T) {
	deviceDate := LoginSessionDeviceDate{IP: "192.168.1.1"}

	// Parse subnets
	_, subnet1, _ := net.ParseCIDR("192.168.1.0/24")
	_, subnet2, _ := net.ParseCIDR("10.0.0.0/8")
	subnets := []*net.IPNet{subnet1, subnet2}

	// IP is within one of the subnets
	inSubnetList, err := deviceDate.IsIPInSubnetList(subnets)
	assert.NoError(t, err)
	assert.True(t, inSubnetList)

	// IP is not within any of the subnets
	_, subnet3, _ := net.ParseCIDR("172.16.0.0/12")
	subnets = []*net.IPNet{subnet3}
	inSubnetList, err = deviceDate.IsIPInSubnetList(subnets)
	assert.NoError(t, err)
	assert.False(t, inSubnetList)
}
