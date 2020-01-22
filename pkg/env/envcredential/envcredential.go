package envcredential

import (
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
)

type credentialResolver struct {
	credManager credential.Manager
}

// NewResolver creates a github resolver instance of Resolver interface
func NewResolver(credManager credential.Manager) *credentialResolver {
	return &credentialResolver{credManager}
}

func (r *credentialResolver) Resolve(name string) (string, error) {
	s := strings.Split(name, "_")
	provider := strings.ToLower(s[1])
	key := strings.ToLower(s[2])
	cred, err := r.credManager.Get(provider)
	if err != nil {
		return "", err
	}
	return cred.Credential[key], nil
}
