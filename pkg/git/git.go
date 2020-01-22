package git

// Credential type that represents git credential
type Credential struct {
	Username string
	Token    string
}

// Options type that represents options for git plain clone operation
type Options struct {
	Credential *Credential
	URL        string
}

//go:generate $GOPATH/bin/moq -out mock_gitrepomanager.go . RepoManager

// RepoManager is an interface that we can use to perform git repository operations
type RepoManager interface {
	PlainClone(path string, o *Options) error
	Pull(path string, o *Options) error
}
