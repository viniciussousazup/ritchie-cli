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

// Session type that represents a session of the user login
type Session struct {
	AccessToken  string `json:"access_token"`
	Organization string `json:"organization"`
	Username     string `json:"username"`
}

type ProviderConfig struct {
	ConfigUrl		string `json:"configUrl"`
	ClientId		string `json:"clientId"`
}

//go:generate $GOPATH/bin/moq -out mock_loginmanager.go . Manager

// Manager is an interface that we can use to perform login operations
type Manager interface {
	Authenticate(organization string) error
	Session() (*Session, error)
}
