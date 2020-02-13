package context

//go:generate $GOPATH/bin/moq -out mock_contextmanager.go . Manager

type Manager interface {
	Set(ctx string) error
	Show() (string, error)
	Delete() error
}
