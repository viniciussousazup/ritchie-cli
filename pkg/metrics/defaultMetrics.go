package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/login"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	urlPatternMetricsUse = "%s/metrics/use"
)

type defaultManager struct {
	serverURL    string
	httpClient   *http.Client
	loginManager login.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(serverUrl string, c *http.Client, l login.Manager) *defaultManager {
	return &defaultManager{serverURL: serverUrl, httpClient: c, loginManager: l}
}

func (d *defaultManager) SendCommand() {
	session, err := d.loginManager.Session()
	if err != nil {
		return
	}

	cmdUse := CmdUse{
		Username: session.Username,
		Cmd: proccessCmd(),
	}

	b, err := json.Marshal(&cmdUse)
	if err != nil {
		return
	}

	url := fmt.Sprintf(urlPatternMetricsUse, d.serverURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", session.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

}

func proccessCmd() string{
	var strs []string

	for i := 0; i < len(os.Args) ; i ++ {
		if i == len(os.Args)-1 {
			strs = append(strs,os.Args[i])
			continue
		}
		strs = append(strs,os.Args[i]+" ")
	}
	return strings.Join(strs, "")

}
