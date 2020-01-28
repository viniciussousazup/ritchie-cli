package git

import (
	"os"

	gogit "gopkg.in/src-d/go-git.v4"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type defaultManager struct{}

// NewDefaultManager creates a default instance of RepoManager interface
func NewDefaultManager() *defaultManager {
	return &defaultManager{}
}

func (*defaultManager) PlainClone(path string, o *Options) error {
	_, err := gogit.PlainClone(path, false, &gogit.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: o.Credential.Username,
			Password: o.Credential.Token,
		},
		URL:      o.URL,
		Progress: os.Stdout,
	})
	return err
}

func (*defaultManager) Pull(path string, o *Options) error {
	r, err := gogit.PlainOpen(path)
	if err != nil {
		return nil
	}
	w, err := r.Worktree()
	if err != nil {
		return nil
	}
	err = w.Pull(&gogit.PullOptions{
		RemoteName: "origin",
		Auth: &githttp.BasicAuth{
			Username: o.Credential.Username,
			Password: o.Credential.Token,
		},
		Progress: os.Stdout,
	})
	return err
}
