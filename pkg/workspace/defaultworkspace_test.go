package workspace

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/login"
	"github.com/ZupIT/ritchie-cli/pkg/tree"
)

var (
	home          string
	homeNotExists string
	serverURL     string
)

func TestMain(m *testing.M) {
	home = fmt.Sprintf("%s/.rit", os.TempDir())
	homeNotExists = fmt.Sprintf("%s/.notexists", os.TempDir())
	serverURL = "https://ritchie-server.itiaws.dev"
	fileutil.CreateIfNotExists(home, 0755)
	os.Exit(m.Run())
}

func TestCheckWorkingDir(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		label string
		in    string
	}{
		{"check working dir", home},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			workman := NewDefaultManager(test.in, serverURL, http.DefaultClient, &tree.ManagerMock{}, &git.RepoManagerMock{}, &credential.ManagerMock{}, &login.ManagerMock{})
			err := workman.CheckWorkingDir()
			is.NoErr(err)
		})
	}
}

func TestInitWorkingDir(t *testing.T) {
	is := is.New(t)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		okResponse := `{"jenkins":[{"field":"username","type":"text"},{"field":"token","type":"password"}],"github":[{"field":"username","type":"text"},{"field":"token","type":"password"}],"gitlab":[{"field":"username","type":"text"},{"field":"token","type":"password"}],"aws":[{"field":"accessKeyId","type":"text"},{"field":"secretAccessKey","type":"password"}],"darwin":[{"field":"username","type":"text"},{"field":"password","type":"password"}]}`
		w.Write([]byte(okResponse))
	})

	workman, teardown := newTestingManager(h)
	defer teardown()

	err := workman.InitWorkingDir()
	is.NoErr(err)
}

func newTestingManager(handler http.Handler) (Manager, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	treeman := &tree.ManagerMock{
		LoadAndSaveTreeFunc: func() error {
			return nil
		},
	}
	repoman := &git.RepoManagerMock{
		PlainCloneFunc: func(path string, o *git.Options) error {
			return nil
		},
		PullFunc: func(path string, o *git.Options) error {
			return nil
		},
	}
	credman := &credential.ManagerMock{
		SaveFunc: func(s *credential.Secret) error {
			return nil
		},
		GetFunc: func(provider string) (*credential.Secret, error) {
			return &credential.Secret{
				Username:   "test",
				Credential: make(map[string]string),
				Provider:   "test",
			}, nil
		},
		ConfigsFunc: func() (credential.Configs, error) {
			return nil, nil
		},
	}

	logman := &login.ManagerMock{
		SessionFunc: func() (*login.Session, error) {
			return &login.Session{}, nil
		},
	}

	workman := NewDefaultManager(home, s.URL, cli, treeman, repoman, credman, logman)

	return workman, s.Close
}
