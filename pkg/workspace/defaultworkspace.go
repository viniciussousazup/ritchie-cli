package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/tree"
)

const (
	urlConfigPattern = "%s/configs"
	usernameKey      = "username"
	tokenKey         = "token"
)

var (
	errAlreadyUpToDate = errors.New("already up-to-date")
)

type formulaConfig struct {
	Provider string `json:"formula_provider"`
	URL      string `json:"formula_url"`
}

// defaultManager is a default implementation of Manager interface
type defaultManager struct {
	ritchieHome    string
	serverURL      string
	httpClient     *http.Client
	treeManager    tree.Manager
	gitRepoManager git.RepoManager
	credManager    credential.Manager
	session        session.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(
	ritchieHome string,
	serverUrl string,
	c *http.Client,
	t tree.Manager,
	g git.RepoManager,
	cr credential.Manager,
	s session.Manager,
) *defaultManager {
	return &defaultManager{
		ritchieHome:    ritchieHome,
		serverURL:      serverUrl,
		httpClient:     c,
		treeManager:    t,
		gitRepoManager: g,
		credManager:    cr,
		session:        s,
	}
}

// CheckWorkingDir default implementation of function Manager.CheckWorkingDir
func (d *defaultManager) CheckWorkingDir() error {
	err := fileutil.CreateIfNotExists(d.ritchieHome, 0755)
	if err != nil {
		return err
	}
	return nil
}

// InitWorkingDir default implementation of function Manager.InitWorkingDir
func (d *defaultManager) InitWorkingDir() error {
	log.Println("Loading user s...")
	s, err := d.session.Get()
	if err != nil {
		return err
	}
	log.Println("done.")

	log.Println("Loading and saving command tree...")
	err = d.treeManager.LoadAndSaveTree()
	if err != nil {
		return err
	}
	log.Println("done.")

	log.Println("Getting formulas...")
	url := fmt.Sprintf(urlConfigPattern, d.serverURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	switch resp.StatusCode {
	case 200:
		frm := &formulaConfig{}
		err := json.NewDecoder(resp.Body).Decode(frm)
		if err != nil {
			return err
		}

		switch frm.Provider {
		case "github":
			err := d.gitFormulaProcess(frm)
			if err != nil {
				return err
			}
		case "s3":
			err := d.s3Process(frm)
			if err != nil {
				return err
			}
		}

	default:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Printf("Status code: %v", resp.StatusCode)
		return errors.New(string(b))
	}

	log.Println("done.")

	return nil
}

func (d *defaultManager) gitFormulaProcess(frm *formulaConfig) error {
	secret, err := d.credManager.Get(frm.Provider)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(formula.PathPattern, d.ritchieHome, "")

	opt := &git.Options{
		Credential: &git.Credential{
			Username: secret.Credential[usernameKey],
			Token:    secret.Credential[tokenKey],
		},
		URL: frm.URL,
	}

	if fileutil.Exists(path) {
		log.Println("Pull formulas...")
		err = d.gitRepoManager.Pull(path, opt)
		if err != nil && err.Error() != errAlreadyUpToDate.Error() {
			return nil
		}
	} else {
		log.Println("Clone formulas...")
		err = d.gitRepoManager.PlainClone(path, opt)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (d *defaultManager) s3Process(frm *formulaConfig) error {
	destPath := fmt.Sprintf(d.ritchieHome+"%s", formula.DirFormula)
	zipFile, err := downloadZipProject(frm.URL, destPath)
	if err != nil {
		return err
	}
	err = unzipFile(zipFile, d.ritchieHome)
	if err != nil {
		return nil
	}
	return nil
}

func downloadZipProject(url, destPath string) (string, error) {
	log.Println("Starting download zip file.")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	file := fmt.Sprintf("%s.zip", destPath)
	out, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	log.Println("Download zip file done.")
	return file, nil
}

func unzipFile(filename, destPath string) error {
	log.Println("Unzip files S3...")

	_ = fileutil.CreateIfNotExists(destPath, 0655)
	err := fileutil.Unzip(filename, destPath)
	if err != nil {
		return err
	}
	err = fileutil.RemoveFile(filename)
	if err != nil {
		return err
	}
	log.Println("Unzip S3 done.")
	return nil
}
