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

func TestNewListAdder(t *testing.T) {

	ritHome := os.TempDir()
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	repoList := NewLister(ritHome, fileManager)
	repoCreator := NewCreator(ritHome, github.NewRepoManager(http.DefaultClient), dirManager, fileManager)
	treeGenerator := tree.NewGenerator(dirManager, fileManager)
	repoAdd := NewAdder(ritHome, repoCreator, treeGenerator, dirManager, fileManager)

	type in struct {
		repoList formula.RepositoryLister
		repoAdd  formula.RepositoryAdder
	}
	tests := []struct {
		name string
		in   in
		want formula.RepositoryAddLister
	}{
		{
			name: "Build with success",
			in: in{
				repoList: repoList,
				repoAdd:  repoAdd,
			},
			want: ListAddManager{
				RepositoryAdder:  repoAdd,
				RepositoryLister: repoList,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewListAdder(tt.in.repoList, tt.in.repoAdd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListAdder() = %v, want %v", got, tt.want)
			}
		})
	}
}
