package context

import (
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type defaultManager struct {
	session session.Manager
}

func NewDefaultManager(s session.Manager) *defaultManager {
	return &defaultManager{session: s}
}

func (d *defaultManager) Set(ctx string) error {
	err := d.session.SetCtx(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d *defaultManager) Show() (string, error) {
	s, err := d.session.Get()
	if err != nil {
		return "", err
	}

	return s.Context, nil
}

func (d *defaultManager) Delete() error {
	s, err := d.session.Get()
	if err != nil {
		return err
	}

	err = d.session.Create(s.AccessToken, s.Username, s.Organization)
	if err != nil {
		return err
	}

	return nil
}
