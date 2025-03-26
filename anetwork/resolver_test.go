package anetwork

import (
	"context"
	"testing"
	"time"
)

// TestGetResolv tests the GetResolv function for reading and parsing the resolv.conf file.
func TestGetResolv(t *testing.T) {
	resolver, err := GetResolv()
	if err != nil {
		t.Errorf("GetResolv() error = %v", err)
	}
	if len(resolver.Nameservers) == 0 {
		t.Error("GetResolv() expected to find at least one nameserver")
	}
}

// TestGetNetResolver tests the GetNetResolver function for creating a custom net.Resolver.
func TestGetNetResolver(t *testing.T) {
	nameserver := "8.8.8.8" // Google's public DNS server
	resolver := GetNetResolver(nameserver)
	if resolver == nil {
		t.Error("GetNetResolver() returned nil")
	}

	// Test resolving using the custom resolver
	_, err := resolver.LookupHost(context.Background(), "example.com")
	if err != nil {
		t.Errorf("GetNetResolver() failed to resolve host: %v", err)
	}
}

// TestGetNetResolver_Default tests the GetNetResolver function with an empty nameserver string.
func TestGetNetResolver_Default(t *testing.T) {
	resolver := GetNetResolver("")
	if resolver == nil {
		t.Error("GetNetResolver() returned nil for default resolver")
	}

	// Test resolving using the default resolver
	_, err := resolver.LookupHost(context.Background(), "example.com")
	if err != nil {
		t.Errorf("GetNetResolver() failed to resolve host using default resolver: %v", err)
	}
}

func TestGetNetResolver_Invalid(t *testing.T) {
	nameserver := "invalid-nameserver"
	resolver := GetNetResolver(nameserver)
	if resolver == nil {
		t.Error("GetNetResolver() returned nil for invalid nameserver")
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Test resolving using the custom resolver with an invalid nameserver
	_, err := resolver.LookupHost(ctx, "example.com")
	if err == nil {
		t.Error("GetNetResolver() expected to fail resolving host with invalid nameserver")
	}
}
