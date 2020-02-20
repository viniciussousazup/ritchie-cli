package workspace

import (
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"

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
	_ = fileutil.CreateIfNotExists(home, 0755)
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
			workman := NewDefaultManager(test.in, &tree.ManagerMock{})
			err := workman.CheckWorkingDir()
			is.NoErr(err)
		})
	}
}
