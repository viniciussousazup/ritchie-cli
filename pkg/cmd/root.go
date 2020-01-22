package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/workspace"
)

type rootCmd struct {
	workspaceManager workspace.Manager
}

// NewRootCmd creates the root for all ritchie commands.
func NewRootCmd(wm workspace.Manager) *cobra.Command {
	o := &rootCmd{wm}

	return &cobra.Command{
		Use:   "rit",
		Short: "rit is a NoOps CLI",
		Long: "A CLI that developers can build and operate\n" +
			"your applications without help from the infra staff.\n" +
			"Complete documentation is available at https://github.com/ZupIT/ritchie-cli",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return o.workspaceManager.CheckWorkingDir()
		},
	}
}
