package aclient_smtp

import (
	"fmt"
	"net/smtp"
)

// loginAuth implements the smtp.Auth interface for LOGIN authentication
type loginAuth struct {
	username, password string
}

// NewLoginAuth creates a new instance of loginAuth with the provided username and password
func NewLoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Start begins the authentication process by returning the mechanism name and initial response
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

// Next continues the authentication process based on the server's challenge
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown from server")
		}
	}
	return nil, nil
}
