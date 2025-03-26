package anetwork

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestToNetIP(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		expected  NetIPs
		expectErr bool
	}{
		{
			name:  "Single IP",
			input: []string{"192.168.1.1"},
			expected: NetIPs{
				{IP: net.ParseIP("192.168.1.1")},
			},
		},
		{
			name:  "Single CIDR",
			input: []string{"10.0.0.0/8"},
			expected: NetIPs{
				{Subnet: &net.IPNet{
					IP:   net.ParseIP("10.0.0.0"),
					Mask: net.CIDRMask(8, 32),
				}},
			},
		},
		{
			name:  "Multiple IPs and CIDRs",
			input: []string{"192.168.1.1", "10.0.0.0/8"},
			expected: NetIPs{
				{IP: net.ParseIP("192.168.1.1")},
				{Subnet: &net.IPNet{
					IP:   net.ParseIP("10.0.0.0"),
					Mask: net.CIDRMask(8, 32),
				}},
			},
		},
		{
			name:      "Invalid IP",
			input:     []string{"invalid-ip"},
			expectErr: true,
		},
		{
			name:      "Mixed valid and invalid IPs",
			input:     []string{"192.168.1.1", "invalid-ip"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToNetIP(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(result))
				for i := range tt.expected {
					if tt.expected[i].IP != nil {
						assert.Equal(t, tt.expected[i].IP.String(), result[i].IP.String())
					} else {
						assert.Equal(t, tt.expected[i].Subnet.String(), result[i].Subnet.String())
					}
				}
			}
		})
	}
}

func TestNetIPs_Contains(t *testing.T) {
	netIPs, err := ToNetIP([]string{"192.168.1.1", "10.0.0.0/8"})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		ipTarget string
		expected bool
	}{
		{
			name:     "Contains single IP",
			ipTarget: "192.168.1.1",
			expected: true,
		},
		{
			name:     "Does not contain single IP",
			ipTarget: "192.168.1.2",
			expected: false,
		},
		{
			name:     "Contains IP in CIDR",
			ipTarget: "10.0.0.1",
			expected: true,
		},
		{
			name:     "Does not contain IP outside CIDR",
			ipTarget: "11.0.0.1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := netIPs.Contains(tt.ipTarget)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestConvertIPToFilename tests the conversion of individual IPs to safe filenames with optional extensions.
func TestConvertIPToFilename(t *testing.T) {
	tests := []struct {
		ip       string
		ext      string
		expected string
	}{
		// Standard IPv4 cases
		{"192.168.1.1", "log", "192_168_1_1.log"},
		{"10.0.0.255", "", "10_0_0_255"}, // No extension
		{"127.0.0.1", ".txt", "127_0_0_1.txt"},
		{"127.0.0.1", "   log   ", "127_0_0_1.log"}, // Trimmed extension
		{"127.0.0.1", ".", "127_0_0_1"},             // Only "." should be ignored

		// Standard IPv6 cases
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "log", "2001_0db8_85a3_0000_0000_8a2e_0370_7334.log"},
		{"fe80::1", "txt", "fe80__1.txt"},
		{"::1", "cfg", "__1.cfg"},
		{"::ffff:192.168.1.1", "log", "__ffff_192_168_1_1.log"},

		// Edge cases
		{"", "log", ".log"}, // Empty input should return ".log"
		{"", "", ""},        // Empty input and no extension should return ""
		{"localhost", "log", "localhost.log"},
	}

	for _, test := range tests {
		t.Run(test.ip, func(t *testing.T) {
			result := ConvertIPToFilename(test.ip, test.ext)
			if result != test.expected {
				t.Errorf("ConvertIPToFilename(%q, %q) = %q; expected %q", test.ip, test.ext, result, test.expected)
			}
		})
	}
}

// TestConvertIPLogsToFilenames tests converting multiple IPs at once with optional extensions.
func TestConvertIPLogsToFilenames(t *testing.T) {
	input := []string{
		"192.168.1.1",
		"10.0.0.255",
		"127.0.0.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"fe80::1",
		"::1",
		"",
		"localhost",
	}

	// Test with ".log" extension
	expectedLogExt := []string{
		"192_168_1_1.log",
		"10_0_0_255.log",
		"127_0_0_1.log",
		"2001_0db8_85a3_0000_0000_8a2e_0370_7334.log",
		"fe80__1.log",
		"__1.log",
		".log",
		"localhost.log",
	}

	result := ConvertIPLogsToFilenames(input, "log")
	for i, res := range result {
		if res != expectedLogExt[i] {
			t.Errorf("ConvertIPLogsToFilenames(log) [%d] = %q; expected %q", i, res, expectedLogExt[i])
		}
	}

	// Test with no extension
	expectedNoExt := []string{
		"192_168_1_1",
		"10_0_0_255",
		"127_0_0_1",
		"2001_0db8_85a3_0000_0000_8a2e_0370_7334",
		"fe80__1",
		"__1",
		"",
		"localhost",
	}

	result = ConvertIPLogsToFilenames(input, "")
	for i, res := range result {
		if res != expectedNoExt[i] {
			t.Errorf("ConvertIPLogsToFilenames(no ext) [%d] = %q; expected %q", i, res, expectedNoExt[i])
		}
	}
}
