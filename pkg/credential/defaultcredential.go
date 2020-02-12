package credential

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	urlCreatePattern = "%s/credentials/%s"
	urlGetPattern    = "%s/credentials/me/%s"
	urlConfigPattern = "%s/credentials/config"
)

type credentialJSON struct {
	Service    string            `json:"service"`
	Username   string            `json:"username"`
	Credential map[string]string `json:"credential"`
}

type defaultManager struct {
	serverURL  string
	httpClient *http.Client
	session    session.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(serverURL string, c *http.Client, s session.Manager) *defaultManager {
	return &defaultManager{serverURL: serverURL, httpClient: c, session: s}
}

func (d *defaultManager) Save(secret *Secret) error {
	s, err := d.session.Get()
	if err != nil {
		return err
	}

	cred := &credentialJSON{
		Service:    secret.Provider,
		Username:   secret.Username,
		Credential: secret.Credential,
	}

	b, err := json.Marshal(&cred)
	if err != nil {
		return err
	}

	path := "me"
	if secret.Username != Me {
		path = "admin"
	} else {
		secret.Username = ""
	}

	url := fmt.Sprintf(urlCreatePattern, d.serverURL, path)
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
		log.Printf("Status code: %v", resp.StatusCode)
		return errors.New(string(b))
	}
}

func (d *defaultManager) Get(provider string) (*Secret, error) {
	s, err := d.session.Get()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlGetPattern, d.serverURL, provider)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-org", s.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		cred := &credentialJSON{}
		json.NewDecoder(resp.Body).Decode(cred)
		sec := &Secret{
			Username:   s.Username,
			Credential: cred.Credential,
			Provider:   provider,
		}
		return sec, nil
	default:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		log.Printf("Status code: %v", resp.StatusCode)
		return nil, errors.New(string(b))
	}
}

func (d *defaultManager) Configs() (Configs, error) {
	s, err := d.session.Get()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlConfigPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-org", s.Organization)
	req.Header.Set("x-ctx", s.Context)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return handler200(resp)
	default:
		return handlerError(resp)
	}
}

func handler200(resp *http.Response) (Configs, error) {
	var cfg Configs
	err := json.NewDecoder(resp.Body).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func handlerError(resp *http.Response) (Configs, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Status code: %v", resp.StatusCode)
	return nil, errors.New(string(b))
}
