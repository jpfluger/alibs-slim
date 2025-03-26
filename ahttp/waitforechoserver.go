package ahttp

import (
	stdContext "context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/alog"
	"io"
	"net"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

// WaitForServerStartEcho waits for an Echo server to start using channels.
// It checks the server's listener address to determine if the server has started.
// The function now accepts a retry interval parameter.
func WaitForServerStartEcho(e *echo.Echo, errChan <-chan error, isTLS bool, retryInterval time.Duration) error {
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var addr net.Addr
			if isTLS {
				addr = e.TLSListenerAddr()
			} else {
				addr = e.ListenerAddr()
			}
			if addr != nil && strings.Contains(addr.String(), ":") {
				alog.LOGGER(alog.LOGGER_APP).Info().Msg("server has started")
				return nil // Server has started
			}
		case err := <-errChan:
			if err == http.ErrServerClosed {
				alog.LOGGER(alog.LOGGER_APP).Info().Msg("server has closed")
				return nil // Server has closed, considered as started for this context
			}
			alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("server start error")
			return err // Return any other error that occurred
		}
	}
}

// WaitForServerStartPing sends HTTP GET requests to the specified URL until the expected response is received or it times out.
// The function now accepts a retry interval parameter and includes logging.
func WaitForServerStartPing(urlRequest string, expectedResults string, pingTimeOutSeconds int, retryInterval time.Duration) error {
	// Validate the URL format.
	if _, err := url2.Parse(urlRequest); err != nil {
		return fmt.Errorf("failed to parse urlRequest: %v", err)
	}

	// Attempt to ping the server until the timeout is reached.
	for ii := 0; ii < pingTimeOutSeconds; ii++ {
		time.Sleep(retryInterval)

		resp, err := http.Get(urlRequest)
		if err != nil {
			alog.LOGGER(alog.LOGGER_APP).Err(err).Str("url", urlRequest).Msg("ping attempt failed")
			continue // Continue pinging if there's an error
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read body of web response for '%s': %v", urlRequest, err)
		}

		bodyString := string(body)
		if expectedResults == "" || bodyString == expectedResults {
			alog.LOGGER(alog.LOGGER_APP).Info().Str("url", urlRequest).Msg("server is ready")
			return nil // Expected result found, server is ready
		}

		alog.LOGGER(alog.LOGGER_APP).Err(err).Str("url", urlRequest).Str("res.body", bodyString).Msg("unexpected server response")
	}

	return fmt.Errorf("server did not start within the timeout period")
}

//package ahttp
//
//import (
//	stdContext "context"
//	"fmt"
//	"github.com/labstack/echo/v4"
//	"io"
//	"net"
//	"net/http"
//	url2 "net/url"
//	"strings"
//	"time"
//)
//
//// WaitForServerStartEcho uses channels waits for an echo server to start.
//func WaitForServerStartEcho(e *echo.Echo, errChan <-chan error, isTLS bool) error {
//	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 200*time.Millisecond)
//	defer cancel()
//
//	ticker := time.NewTicker(5 * time.Millisecond)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ctx.Done():
//			return ctx.Err()
//		case <-ticker.C:
//			var addr net.Addr
//			if isTLS {
//				addr = e.TLSListenerAddr()
//			} else {
//				addr = e.ListenerAddr()
//			}
//			if addr != nil && strings.Contains(addr.String(), ":") {
//				return nil // was started
//			}
//		case err := <-errChan:
//			if err == http.ErrServerClosed {
//				return nil
//			}
//			return err
//		}
//	}
//}
//
//func WaitForServerStartPing(urlRequest string, expectedResults string, pingTimeOutSeconds int) error {
//	if _, err := url2.Parse(urlRequest); err != nil {
//		return fmt.Errorf("failed to parse urlRequest; %v", err)
//	}
//
//	ii := 0
//	var errPing error
//	for ii < pingTimeOutSeconds {
//		time.Sleep(1 * time.Second)
//		ii++
//
//		resp, err := http.Get(urlRequest)
//		if err != nil {
//			errPing = fmt.Errorf("could not ping url %s; %v", urlRequest, err)
//			continue
//		}
//		defer resp.Body.Close()
//
//		body, err := io.ReadAll(resp.Body)
//		if err != nil {
//			errPing = fmt.Errorf("could not read body of web response for '%s'; %v", urlRequest, err)
//			break
//		}
//
//		bodyString := string(body)
//		if expectedResults == "" || bodyString == expectedResults {
//			errPing = nil
//			break
//		}
//
//		errPing = fmt.Errorf("body string does not match expected results for '%s'", urlRequest)
//		break
//	}
//
//	return errPing
//}
