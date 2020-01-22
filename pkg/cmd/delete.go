package cmd

import "github.com/spf13/cobra"

func NewDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete SUBCOMMAND",
		Short: "Delete objects",
		Long:  `Delete objects like users, etc.`,
	}
}
