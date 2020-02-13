package session

const (
	sessionFilePattern    = "%s/.session"
	passphraseFilePattern = "%s/.passphrase"
)

type Session struct {
	AccessToken  string `json:"access_token"`
	Organization string `json:"organization"`
	Username     string `json:"username"`
	Context      string `json:"context"`
}

//go:generate $GOPATH/bin/moq -out mock_sessionmanager.go . Manager

type Manager interface {
	Create(token, username, organization string) error
	Get() (*Session, error)
	SetCtx(ctx string) error
}
