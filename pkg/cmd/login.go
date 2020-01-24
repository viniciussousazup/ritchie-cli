package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/login"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// LoginCmd type for init command
type LoginCmd struct {
	loginManager login.Manager
}

// NewLoginCmd creates new cmd instance
func NewLoginCmd(loginManager login.Manager) *cobra.Command {
	o := &LoginCmd{loginManager}
	return &cobra.Command{
		Use:   "login",
		Short: "User login",
		Long:  "Authenticates the user in an organization to which he belongs",
		RunE: func(cmd *cobra.Command, args []string) error {
			org, err := prompt.String("Login [Organization]: ", true)
			if err != nil {
				return err
			}
			log.Println("Starting login...")
			err = o.loginManager.Authenticate(org)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
