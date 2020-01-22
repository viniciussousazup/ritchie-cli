package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/workspace"
)

// InitCmd type for init command
type InitCmd struct {
	workspaceManager workspace.Manager
}

// NewInitCmd creates new cmd instance
func NewInitCmd(wm workspace.Manager) *cobra.Command {
	o := &InitCmd{wm}
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a working directory",
		Long:  `Initialize a working directory ($USER_HOME/.rit) containing RIT configuration files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.workspaceManager.InitWorkingDir()
			if err != nil {
				return err
			}
			return nil
		},
	}
}
