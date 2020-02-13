package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/context"
)

type setContextCmd struct {
	ctxManager context.Manager
}

func NewSetContextCmd(ctxManager context.Manager) *cobra.Command {
	c := setContextCmd{ctxManager: ctxManager}

	return &cobra.Command{
		Use:     "context",
		Short:   "Set context for Ritchie-cli",
		Example: "rit set context my_context_name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.prompt(args)
		},
	}
}

func (c *setContextCmd) prompt(args []string) error {

	if len(args) < 1 {
		return errors.New("set a name for the context is mandatory")
	}

	err := c.ctxManager.Set(args[0])
	if err != nil {
		return err
	}

	fmt.Println("Set context successful!")
	return nil
}
