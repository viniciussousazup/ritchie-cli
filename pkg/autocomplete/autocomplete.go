package autocomplete

//go:generate $GOPATH/bin/moq -out mock_autocompletemanager.go . Manager

type Manager interface {
	Handle(shellName string) (string, error)
}