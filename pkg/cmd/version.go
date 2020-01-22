package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	BuildDate = "unknown"
)

// NewVersionCmd creates a new cmd instance
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Long:  `Display version and build information about rit.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rit %s\n", Version)
			fmt.Printf("  Build date: %s\n", BuildDate)
			fmt.Printf("  Built with: %s\n", runtime.Version())
		},
	}
}
