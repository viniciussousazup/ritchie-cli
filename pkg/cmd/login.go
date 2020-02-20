package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/login"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// loginCmd type for init command
type loginCmd struct {
	loginManager login.Manager
}

// NewLoginCmd creates new cmd instance
func NewLoginCmd(l login.Manager) *cobra.Command {
	o := &loginCmd{l}
	return &cobra.Command{
		Use:   "login",
		Short: "User login",
		Long:  "Authenticates the user in an organization to which he belongs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt()
		},
	}
}

func (o *loginCmd) prompt() error {
	org, err := prompt.String("Login [Organization]: ", true)
	if err != nil {
		return err
	}

	log.Println("Starting login...")
	if err := o.loginManager.Authenticate(org, Version); err != nil {
		return err
	}
	return nil
}
