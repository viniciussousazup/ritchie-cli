package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/context"
)

type showContextCmd struct {
	ctxManager context.Manager
}

func NewShowContextCmd(ctxManager context.Manager) *cobra.Command {
	c := showContextCmd{ctxManager: ctxManager}

	return &cobra.Command{
		Use:     "context",
		Short:   "Show current context",
		Example: "rit show context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.prompt()
		},
	}
}

func (c *showContextCmd) prompt() error {
	ctx, err := c.ctxManager.Show()
	if err != nil {
		return err
	}

	if ctx == "" {
		fmt.Println("You don't have a defined context")
		return nil
	}

	fmt.Println(fmt.Sprintf("Current context: %s", ctx))

	return nil
}
