package credential

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/login"
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
	serverURL    string
	httpClient   *http.Client
	loginManager login.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(serverURL string, c *http.Client, l login.Manager) *defaultManager {
	return &defaultManager{serverURL: serverURL, httpClient: c, loginManager: l}
}

func (d *defaultManager) Save(secret *Secret) error {
	session, err := d.loginManager.Session()
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
	req.Header.Set("x-org", session.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
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
	session, err := d.loginManager.Session()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlGetPattern, d.serverURL, provider)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-org", session.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
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
			Username:   session.Username,
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
	session, err := d.loginManager.Session()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlConfigPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-org", session.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		var cfg Configs
		_ = json.NewDecoder(resp.Body).Decode(&cfg)
		return cfg, nil
	default:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		log.Printf("Status code: %v", resp.StatusCode)
		return nil, errors.New(string(b))
	}
}
