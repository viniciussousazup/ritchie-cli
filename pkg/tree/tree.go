package tree

// Representation type
type Representation struct {
	Commands []Command `json:"commands"`
	Version  string    `json:"version"`
}

// Command type
type Command struct {
	Parent  string  `json:"parent"`
	Usage   string  `json:"usage"`
	Help    string  `json:"help"`
	Formula Formula `json:"formula"`
}

// Formula type
type Formula struct {
	Path    string `json:"path"`
	Bin     string `json:"bin"`
	Config  string `json:"config"`
	RepoUrl string `json:"repoUrl"`
}

//go:generate $GOPATH/bin/moq -out mock_treemanager.go . Manager

// Manager is an interface that we can use to perform tree operations
type Manager interface {
	GetLocalTree() (*Representation, error)
	LoadAndSaveTree() error
}
