package credential

const (
	// Me credential
	Me = "me"
	// Admin credential
	Admin = "admin"
)

// Secret type that represents a credential secret
type Secret struct {
	Username   string
	Credential map[string]string
	Provider   string
}

// Config type that represents a credential config from server
type Config struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

type Configs map[string][]Config

//go:generate $GOPATH/bin/moq -out mock_credentialmanager.go . Manager

// Manager is an interface that we can use to perform git credential operations
type Manager interface {
	Configs() (Configs, error)
	Save(s *Secret) error
	Get(provider string) (*Secret, error)
}
