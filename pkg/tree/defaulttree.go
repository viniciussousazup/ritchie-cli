package tree

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/login"
)

const (
	treeCmdPattern = "%s/.cmd_tree.json"
	urlGetPattern  = "%s/tree"
)

// defaultManager is a default implementation of Manager interface
type defaultManager struct {
	ritchieHome  string
	serverURL    string
	httpClient   *http.Client
	loginManager login.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(ritchieHome, serverURL string, c *http.Client, l login.Manager) *defaultManager {
	return &defaultManager{ritchieHome: ritchieHome, serverURL: serverURL, httpClient: c, loginManager: l}
}

// GetLocalTree default implementation of function Manager.GetLocalTree
func (d *defaultManager) GetLocalTree() (*Representation, error) {
	treeCmdFile := fmt.Sprintf(treeCmdPattern, d.ritchieHome)
	if !fileutil.Exists(treeCmdFile) {
		return nil, nil
	}
	treeFile, err := ioutil.ReadFile(treeCmdFile)
	if err != nil {
		return nil, err
	}

	treeCmd := &Representation{}
	err = json.Unmarshal(treeFile, treeCmd)
	if err != nil {
		return nil, err
	}

	return treeCmd, nil
}

// LoadAndSaveTree default implementation of function Manager.SaveTree
func (d *defaultManager) LoadAndSaveTree() error {
	session, err := d.loginManager.Session()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlGetPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-org", session.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		treeCmdFile := fmt.Sprintf(treeCmdPattern, d.ritchieHome)
		err = ioutil.WriteFile(treeCmdFile, body, 0644)
		if err != nil {
			return err
		}
	default:
		log.Printf("Status code: %v", resp.StatusCode)
		return errors.New(string(body))
	}

	return nil
}
