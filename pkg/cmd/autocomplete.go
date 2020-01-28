package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/spf13/cobra"
)

const (
	zsh  = "zsh"
	bash = "bash"
)

// autocompleteCmd type for set autocomplete command
type autocompleteCmd struct {
	manager autocomplete.Manager
}

// NewAutocompleteCmd creates a new cmd instance
func NewAutocompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "completion SUBCOMMAND",
		Short:   "Add autocomplete for terminal",
		Long:    `Add autocomplete for terminal, Available for (bash, zsh).`,
		Example: "rit completion zsh",
	}
}

// NewAutocompleteZsh creates a new cmd instance zsh
func NewAutocompleteZsh(m autocomplete.Manager) *cobra.Command {
	o := &autocompleteCmd{m}

	return &cobra.Command{
		Use:     zsh,
		Short:   "Add zsh autocomplete for terminal",
		Long:    "Add zsh autocomplete for terminal",
		Example: "rit completion zsh",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt(zsh)
		},
	}
}

// NewAutocompleteBash creates a new cmd instance zsh
func NewAutocompleteBash(m autocomplete.Manager) *cobra.Command {
	o := &autocompleteCmd{m}

	return &cobra.Command{
		Use:     bash,
		Short:   "Add bash autocomplete for terminal",
		Long:    "Add bash autocomplete for terminal",
		Example: "rit completion bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt(bash)
		},
	}
}

func (s *autocompleteCmd) prompt(shellName string) error {
	c, err := s.manager.Handle(shellName)
	if err != nil {
		return err
	}

	fmt.Println(c)
	return nil
}
