package ahttp

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestWaitForServerStartEcho tests the WaitForServerStartEcho function.
func TestWaitForServerStartEcho(t *testing.T) {
	e := echo.New()
	errChan := make(chan error, 1)
	go func() {
		errChan <- errors.New("test error")
	}()

	err := WaitForServerStartEcho(e, errChan, false, 5*time.Millisecond)
	if err == nil || err.Error() != "test error" {
		t.Errorf("WaitForServerStartEcho() should return the error sent to the channel")
	}

	// Simulate server start
	listener, _ := net.Listen("tcp", ":0")
	e.Listener = listener
	defer listener.Close()

	err = WaitForServerStartEcho(e, errChan, false, 5*time.Millisecond)
	if err != nil {
		t.Errorf("WaitForServerStartEcho() should return nil when server starts, got: %v", err)
	}
}

// TestWaitForServerStartPing tests the WaitForServerStartPing function.
func TestWaitForServerStartPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	err := WaitForServerStartPing(server.URL, "OK", 5, 10*time.Millisecond)
	if err != nil {
		t.Errorf("WaitForServerStartPing() should return nil when server responds with expected result, got: %v", err)
	}

	err = WaitForServerStartPing("http://invalid-url", "OK", 1, 10*time.Millisecond)
	if err == nil {
		t.Errorf("WaitForServerStartPing() should return an error for invalid URL")
	}
}

// Mock functions and types for testing
type mockListener struct {
	net.Listener
	addr net.Addr
}

func (m *mockListener) Addr() net.Addr {
	return m.addr
}

type mockAddr struct {
	network string
	str     string
}

func (m *mockAddr) Network() string {
	return m.network
}

func (m *mockAddr) String() string {
	return m.str
}

func TestMain(m *testing.M) {
	// Mock echo server listener address
	e := echo.New()
	e.Listener = &mockListener{addr: &mockAddr{network: "tcp", str: "127.0.0.1:0"}}

	// Run the tests
	m.Run()
}
