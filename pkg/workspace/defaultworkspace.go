package workspace

import (
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/tree"
)

// defaultManager is a default implementation of Manager interface
type defaultManager struct {
	ritchieHome string
	treeManager tree.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(ritchieHome string, t tree.Manager) *defaultManager {
	return &defaultManager{ritchieHome: ritchieHome, treeManager: t}
}

// CheckWorkingDir default implementation of function Manager.CheckWorkingDir
func (d *defaultManager) CheckWorkingDir() error {
	err := fileutil.CreateIfNotExists(d.ritchieHome, 0755)
	if err != nil {
		return err
	}
	return nil
}
