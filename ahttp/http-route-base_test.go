package ahttp

import (
	"github.com/jpfluger/alibs-slim/asessions"
	"sync"
	"testing"
)

// TestRouteBase_GetRouteId tests the GetRouteId method of RouteBase.
func TestRouteBase_GetRouteId(t *testing.T) {
	rb := RouteBase{
		RouteId: "test-route",
	}
	if rb.GetRouteId() != "test-route" {
		t.Errorf("GetRouteId() = %v, want %v", rb.GetRouteId(), "test-route")
	}
}

// TestRouteBase_GetMethod tests the GetMethod method of RouteBase.
func TestRouteBase_GetMethod(t *testing.T) {
	rb := RouteBase{
		Method: "GET",
	}
	if rb.GetMethod() != "GET" {
		t.Errorf("GetMethod() = %v, want %v", rb.GetMethod(), "GET")
	}
}

// TestRouteBase_GetPath tests the GetPath method of RouteBase.
func TestRouteBase_GetPath(t *testing.T) {
	rb := RouteBase{
		Path: "/test/path",
	}
	if rb.GetPath() != "/test/path" {
		t.Errorf("GetPath() = %v, want %v", rb.GetPath(), "/test/path")
	}
}

// TestRouteBase_GetPerms tests the GetPerms method of RouteBase.
func TestRouteBase_GetPerms(t *testing.T) {
	rb := RouteBase{
		Perms: asessions.MustNewPermSetByPair("perm1", "R"),
	}
	if rb.GetPerms().HasPermS("perm1:R") != true {
		t.Errorf("GetPerms()[\"perm1\"] = %v, want %v", rb.GetPerms()["perm1"], true)
	}
}

// TestRouteBase_GetRouteNotFoundId tests the GetRouteNotFoundId method of RouteBase.
func TestRouteBase_GetRouteNotFoundId(t *testing.T) {
	rb := RouteBase{
		RouteNotFoundId: "not-found-route",
	}
	if rb.GetRouteNotFoundId() != "not-found-route" {
		t.Errorf("GetRouteNotFoundId() = %v, want %v", rb.GetRouteNotFoundId(), "not-found-route")
	}
}

// TestRouteBase_ConcurrentAccess tests that concurrent access to RouteBase methods does not cause race conditions.
func TestRouteBase_ConcurrentAccess(t *testing.T) {
	rb := RouteBase{
		RouteId:         "concurrent-route",
		Method:          "POST",
		Path:            "/concurrent/path",
		Perms:           asessions.MustNewPermSetByPair("perm2", "R"),
		RouteNotFoundId: "concurrent-not-found",
	}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		if rb.GetRouteId() != "concurrent-route" {
			t.Errorf("Concurrent GetRouteId() failed")
		}
	}()

	go func() {
		defer wg.Done()
		if rb.GetMethod() != "POST" {
			t.Errorf("Concurrent GetMethod() failed")
		}
	}()

	go func() {
		defer wg.Done()
		if rb.GetPath() != "/concurrent/path" {
			t.Errorf("Concurrent GetPath() failed")
		}
	}()

	go func() {
		defer wg.Done()
		if rb.GetPerms().HasPermS("perm2:R") != true {
			t.Errorf("Concurrent GetPerms() failed")
		}
	}()

	go func() {
		defer wg.Done()
		if rb.GetRouteNotFoundId() != "concurrent-not-found" {
			t.Errorf("Concurrent GetRouteNotFoundId() failed")
		}
	}()

	wg.Wait()
}
