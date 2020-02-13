package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/context"
)

type deleteContextCmd struct {
	ctxManager context.Manager
}

func NewDeleteContextCmd(ctxManager context.Manager) *cobra.Command {
	c := deleteContextCmd{ctxManager: ctxManager}

	return &cobra.Command{
		Use:     "context",
		Short:   "Delete context for Ritchie-cli",
		Example: "rit delete context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.prompt()
		},
	}
}

func (c *deleteContextCmd) prompt() error {
	err := c.ctxManager.Delete()
	if err != nil {
		return err
	}

	fmt.Println("Delete context successful!")
	return nil
}
