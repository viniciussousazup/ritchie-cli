package login

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"io"
	"net/http"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
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
func NewDefaultManager(homePath, serverURL string, httpClient *http.Client) *defaultManager {
	return &defaultManager{homePath, serverURL, httpClient}
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
		session := &Session{}
		json.NewDecoder(resp.Body).Decode(session)
		session.Username = cred.Username
		session.Organization = cred.Organization
		b, err := json.Marshal(session)
		if err != nil {
			return err
		}
		id, err := machineid.ID()
		if err != nil {
			return err
		}
		h := md5.New()
		io.WriteString(h, passphrase)
		io.WriteString(h, id)
		cipher := cryptoutil.Encrypt(string(h.Sum(nil)), string(b))
		sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)
		err = fileutil.WriteFile(sessFilePath, []byte(cipher))
		if err != nil {
			return err
		}
		err = os.Chmod(sessFilePath, 0600)
		if err != nil {
			return err
		}
		return nil
	case 401:
		return ErrBadCredential
	case 503:
		return ErrServiceUnavailable
	default:
		return ErrUnknown
	}
}

func (d *defaultManager) Session() (*Session, error) {
	sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)
	if !fileutil.Exists(sessFilePath) {
		return nil, errors.New("Please, you need to login first")
	}
	b, err := fileutil.ReadFile(sessFilePath)
	if err != nil {
		return nil, err
	}
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	h := md5.New()
	io.WriteString(h, passphrase)
	io.WriteString(h, id)
	plain := cryptoutil.Decrypt(string(h.Sum(nil)), string(b))
	session := &Session{}
	if err := json.Unmarshal([]byte(plain), session); err != nil {
		return nil, err
	}
	return session, nil
}
