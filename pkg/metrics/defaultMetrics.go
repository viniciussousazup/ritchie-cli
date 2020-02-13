package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	urlPattern = "%s/metrics/use"
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

func (d *defaultManager) SendCommand() {
	s, err := d.session.Get()
	if err != nil {
		return
	}

	cmdUse := CmdUse{
		Username: s.Username,
		Cmd:      proccessCmd(),
	}

	b, err := json.Marshal(&cmdUse)
	if err != nil {
		return
	}

	url := fmt.Sprintf(urlPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-org", s.Organization)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
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

func proccessCmd() string {
	var strs []string

	for i := 0; i < len(os.Args); i++ {
		if i == len(os.Args)-1 {
			strs = append(strs, os.Args[i])
			continue
		}
		strs = append(strs, os.Args[i]+" ")
	}
	return strings.Join(strs, "")

}
