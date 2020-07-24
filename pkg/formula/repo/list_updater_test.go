package repo

import (
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewListUpdater(t *testing.T) {

	ritHome := os.TempDir()
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	repoList := NewLister(ritHome, fileManager)
	repoCreator := NewCreator(ritHome, github.NewRepoManager(http.DefaultClient), dirManager, fileManager)
	repoListCreator := NewListCreator(repoList, repoCreator)
	treeGenerator := tree.NewGenerator(dirManager, fileManager)
	repoUpdate := NewUpdater(ritHome, repoListCreator, treeGenerator, fileManager)

	type in struct {
		repoList   formula.RepositoryLister
		repoUpdate formula.RepositoryUpdater
	}
	tests := []struct {
		name string
		in   in
		want formula.RepositoryListUpdater
	}{
		{
			name: "Build with success",
			in: in{
				repoList:   repoList,
				repoUpdate: repoUpdate,
			},
			want: ListUpdateManager{
				RepositoryLister:  repoList,
				RepositoryUpdater: repoUpdate,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewListUpdater(tt.in.repoList, tt.in.repoUpdate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListUpdater() = %v, want %v", got, tt.want)
			}
		})
	}
}
