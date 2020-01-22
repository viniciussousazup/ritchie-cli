package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/spf13/cobra"
)

// AutocompleteCmd type for set autocomplete command
type AutocompleteCmd struct {
	manager autocomplete.Manager
}

// NewAutocompleteCmd creates a new cmd instance
func NewAutocompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion SUBCOMMAND",
		Short: "Add autocomplete for terminal",
		Long:  `Add autocomplete for terminal, Available for (bash, zsh).`,
		Example: "rit completion zsh",
	}
}

// NewAutocompleteZsh creates a new cmd instance zsh
func NewAutocompleteZsh(m autocomplete.Manager) *cobra.Command {
	o := &AutocompleteCmd{m}

	return &cobra.Command{
		Use:   "zsh",
		Short: "Add zsh autocomplete for terminal",
		Long:  "Add zsh autocomplete for terminal",
		Example: "rit completion zsh",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt("zsh")
		},
	}
}

// NewAutocompleteBash creates a new cmd instance zsh
func NewAutocompleteBash(m autocomplete.Manager) *cobra.Command {
	o := &AutocompleteCmd{m}

	return &cobra.Command{
		Use:   "bash",
		Short: "Add bash autocomplete for terminal",
		Long:  "Add bash autocomplete for terminal",
		Example: "rit completion bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt("bash")
		},
	}
}

func (s *AutocompleteCmd) prompt(shellName string) error {
	complete, err := s.manager.Handle(shellName)
	if err != nil {
		return err
	}

	fmt.Println(complete)
	return nil
}
