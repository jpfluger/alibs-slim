package aclient_http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func DoHTTPGet(fullURL string, hob *HOB) ([]byte, string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered panic in DoHTTPGet: %v\n", r)
		}
	}()

	if hob == nil {
		return nil, "", fmt.Errorf("HOB is nil")
	}

	client := &http.Client{Timeout: time.Duration(hob.ConnectionTimeout) * time.Second}

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	respContentType := resp.Header.Get("Content-Type")
	if hob.ExpectedType != "" && !strings.HasPrefix(respContentType, hob.ExpectedType) {
		return nil, respContentType, fmt.Errorf("unexpected content type: %s", respContentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, respContentType, err
	}

	return body, respContentType, nil
}

func DoHTTPPost(fullURL string, hob *HOB) ([]byte, string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered panic in DoHTTPPost: %v\n", r)
		}
	}()

	if hob == nil {
		return nil, "", fmt.Errorf("HOB is nil")
	}

	hob.ContentType = strings.TrimSpace(hob.ContentType)
	if hob.ContentType == "" {
		return nil, "", fmt.Errorf("content type not defined")
	}

	var data []byte
	var err error

	if len(hob.Raw) > 0 {
		data = hob.Raw
	} else {
		switch hob.ContentType {
		case "application/json":
			data, err = json.Marshal(hob.Payload)
		case "application/xml":
			data, err = xml.Marshal(hob.Payload)
		default:
			return nil, "", fmt.Errorf("unsupported content type: %s", hob.ContentType)
		}
		if err != nil {
			return nil, "", err
		}
	}

	client := &http.Client{Timeout: time.Duration(hob.ConnectionTimeout) * time.Second}
	resp, err := client.Post(fullURL, hob.ContentType, bytes.NewBuffer(data))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	respContentType := resp.Header.Get("Content-Type")
	if hob.ExpectedType != "" && !strings.HasPrefix(respContentType, hob.ExpectedType) {
		return nil, respContentType, fmt.Errorf("unexpected content type: %s", respContentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, respContentType, err
	}

	return body, respContentType, nil
}
