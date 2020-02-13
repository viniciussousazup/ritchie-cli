package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	urlPattern = "%s/users"
)

type defaultManager struct {
	serverURL  string
	httpClient *http.Client
	session    session.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(serverUrl string, c *http.Client, s session.Manager) *defaultManager {
	return &defaultManager{serverURL: serverUrl, httpClient: c, session: s}
}

func (d *defaultManager) Create(user *Definition) error {
	s, err := d.session.Get()
	if err != nil {
		return err
	}

	b, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", s.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 201:
		return nil
	default:
		return errors.New(string(b))
	}
}

func (d *defaultManager) Delete(user *Definition) error {
	s, err := d.session.Get()
	if err != nil {
		return err
	}

	b, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", s.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	res, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case 200:
		return nil
	case 404:
		return errors.New("user not found")
	default:
		return errors.New(string(b))
	}
}
