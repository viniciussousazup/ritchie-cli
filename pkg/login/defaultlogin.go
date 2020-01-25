package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/crypto"
	"github.com/ZupIT/ritchie-cli/pkg/file"
	"github.com/denisbrodbeck/machineid"
	"net/http"
	"os"
)

const (
	urlPattern         = "%s/login"
	sessionFilePattern = "%s/.session"

	// AES passphrase
	passphrase = "zYtBIK67fCmhrU0iUbPQ1Cf9"
)

type defaultManager struct {
	homePath   string
	serverURL  string
	httpClient *http.Client
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(homePath, serverURL string, c *http.Client) *defaultManager {
	return &defaultManager{homePath, serverURL, c}
}

func (d *defaultManager) Authenticate(cred *Credential) error {
	b, err := json.Marshal(&cred)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", cred.Organization)
	resp, err := d.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return d.handler200(resp, cred)
	case 401:
		return ErrBadCredential
	case 503:
		return ErrServiceUnavailable
	default:
		return ErrUnknown
	}
}

func (d *defaultManager) handler200(resp *http.Response, cred *Credential) error {
	s := &Session{}
	err := json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		return err
	}

	s.Username = cred.Username
	s.Organization = cred.Organization
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	id, err := machineid.ID()
	if err != nil {
		return err
	}

	cipher := crypto.Encrypt(passphrase, id, string(b))
	sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)

	err = file.WriteFile(sessFilePath, []byte(cipher))
	if err != nil {
		return err
	}

	err = os.Chmod(sessFilePath, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (d *defaultManager) Session() (*Session, error) {
	sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)
	if !file.Exists(sessFilePath) {
		return nil, errors.New("Please, you need to login first")
	}
	b, err := file.ReadFile(sessFilePath)
	if err != nil {
		return nil, err
	}
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	plain := crypto.Decrypt(passphrase, id, string(b))
	session := &Session{}
	if err := json.Unmarshal([]byte(plain), session); err != nil {
		return nil, err
	}
	return session, nil
}
