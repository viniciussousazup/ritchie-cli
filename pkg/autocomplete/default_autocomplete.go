package autocomplete

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

const pathUrl = "%s/auto-complete/%s"

type defaultManager struct {
	serverURL    string
	httpClient   *http.Client
}

func NewDefaultManager(serverURL string, c *http.Client) *defaultManager {
	return &defaultManager{serverURL: serverURL, httpClient: c}
}

func (d *defaultManager) Handle(shellName string) (string, error) {

	if !sliceutil.Contains(supportedAutocomplete(), shellName) {
		return "", errors.New("autocomplete for this terminal is not supported")
	}

	url := fmt.Sprintf(pathUrl, d.serverURL, shellName)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func supportedAutocomplete() []string {
	return []string{"bash", "zsh"}
}
