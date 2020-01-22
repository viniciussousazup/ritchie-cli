package login

import (
	"errors"
)

var (
	// ErrBadCredential bad credential error
	ErrBadCredential = errors.New("Bad credentials")

	// ErrServiceUnavailable service unavailable error
	ErrServiceUnavailable = errors.New("Login service unavailable")

	// ErrUnknown unknown error
	ErrUnknown = errors.New("Unknown error. Please, try again")
)

// Credential type that represents a credential of the organization
type Credential struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
}

// Session type that represents a session of the user login
type Session struct {
	AccessToken  string `json:"access_token"`
	Organization string `json:"organization"`
	Username     string `json:"username"`
}

//go:generate $GOPATH/bin/moq -out mock_loginmanager.go . Manager

// Manager is an interface that we can use to perform login operations
type Manager interface {
	Authenticate(cred *Credential) error
	Session() (*Session, error)
}
