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


type ProviderConfig struct {
	Url      string `json:"url"`
	ClientId string `json:"clientId"`
}

//go:generate $GOPATH/bin/moq -out mock_loginmanager.go . Manager

// Manager is an interface that we can use to perform login operations
type Manager interface {
	Authenticate(organization, version string) error
}
