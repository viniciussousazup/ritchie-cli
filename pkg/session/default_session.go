package session

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"

	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type DefaultManager struct {
	homePath string
}

func NewDefaultManager(homePath string) *DefaultManager {
	return &DefaultManager{homePath: homePath}
}

func (d *DefaultManager) Create(token, username, organization string) error {
	session := &Session{
		AccessToken:  token,
		Organization: organization,
		Username:     username,
	}

	err := d.encrypt(session)
	if err != nil {
		return err
	}

	return nil
}

func (d *DefaultManager) Get() (*Session, error) {
	sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)
	if !fileutil.Exists(sessFilePath) {
		return nil, errors.New("please, you need to login first")
	}

	b, err := fileutil.ReadFile(sessFilePath)
	if err != nil {
		return nil, err
	}

	h, err := d.md5()
	if err != nil {
		return nil, err
	}

	plain := cryptoutil.Decrypt(string(h.Sum(nil)), string(b))
	session := &Session{}

	if err := json.Unmarshal([]byte(plain), session); err != nil {
		return nil, err
	}
	return session, nil
}

func (d *DefaultManager) SetCtx(ctx string) error {
	session, err := d.Get()
	if err != nil {
		return err
	}

	session.Context = ctx

	err = d.encrypt(session)
	if err != nil {
		return err
	}

	return nil
}

func (d *DefaultManager) encrypt(session *Session) error {
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}

	h, err := d.md5()
	if err != nil {
		return err
	}

	cipher := cryptoutil.Encrypt(string(h.Sum(nil)), string(b))

	err = d.writeFile(cipher)
	if err != nil {
		return err
	}

	return nil
}

func (d *DefaultManager) md5() (hash.Hash, error) {
	passphrase, err := d.readPassPhrase()
	if err != nil {
		return nil, err
	}

	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}

	h := md5.New()
	_, _ = io.WriteString(h, passphrase)
	_, _ = io.WriteString(h, id)

	return h, nil
}

func (d *DefaultManager) writeFile(cipher string) error {
	sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)

	err := fileutil.WriteFile(sessFilePath, []byte(cipher))
	if err != nil {
		return err
	}

	err = os.Chmod(sessFilePath, 0600)
	if err != nil {
		return nil
	}
	return nil
}

func (d *DefaultManager) readPassPhrase() (string, error) {
	passPhraseFilePath := fmt.Sprintf(passphraseFilePattern, d.homePath)
	if !fileutil.Exists(passPhraseFilePath) {
		passPhrase := uuid.New().String()
		err := fileutil.WriteFile(passPhraseFilePath, []byte(passPhrase))
		if err != nil {
			return "", err
		}
	}
	p, err := fileutil.ReadFile(passPhraseFilePath)
	if err != nil {
		return "", err
	}
	return string(p), nil
}
