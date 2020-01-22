package workspace

import (
	"errors"
)

var (
	// ErrWorkingDirNotInitiliazed represents an error while the working directory not exists yet
	ErrWorkingDirNotInitiliazed = errors.New("Please, execute [rit init] command to initialize the working directory")
)

//go:generate $GOPATH/bin/moq -out mock_workspacemanager.go . Manager

// Manager is an interface that we can use to perform workspace operations
type Manager interface {
	// WorkingDir checks workspace setup
	CheckWorkingDir() error
	// InitWorkingDir creates working dir
	InitWorkingDir() error
}
