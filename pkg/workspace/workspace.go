package workspace

//go:generate $GOPATH/bin/moq -out mock_workspacemanager.go . Manager

// Manager is an interface that we can use to perform workspace operations
type Manager interface {
	// WorkingDir checks workspace setup
	CheckWorkingDir() error
	// InitWorkingDir creates working dir
	InitWorkingDir() error
}
