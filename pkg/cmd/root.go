package cmd

import (
	"fmt"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/workspace"
	"github.com/spf13/cobra"
)

const versionMsg = "%s\n  Build date: %s\n  Built with: %s\n"

var (
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	BuildDate = "unknown"
)

type rootCmd struct {
	workspaceManager workspace.Manager
}

// NewRootCmd creates the root for all ritchie commands.
func NewRootCmd(wm workspace.Manager) *cobra.Command {
	o := &rootCmd{wm}

	return &cobra.Command{
		Use:     "rit",
		Version: o.ritVersion(),
		Short:   "rit is a NoOps CLI",
		Long: `A CLI that developers can build and operate\n
			your applications without help from the infra staff.\n
			Complete documentation is available at https://github.com/ZupIT/ritchie-cli`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return o.workspaceManager.CheckWorkingDir()
		},
	}
}

func (o *rootCmd) ritVersion() string {
	return fmt.Sprintf(versionMsg, Version, BuildDate, runtime.Version())
}
