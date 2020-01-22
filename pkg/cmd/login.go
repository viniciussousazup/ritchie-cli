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
			user, err := prompt.String("Login [Username]: ", true)
			if err != nil {
				return err
			}
			passw, err := prompt.Password("Login [Password]: ")
			if err != nil {
				return err
			}

			cred := &login.Credential{
				Username:     user,
				Password:     passw,
				Organization: org,
			}

			err = o.loginManager.Authenticate(cred)
			if err != nil {
				return err
			}

			log.Println("Login successful!")
			return nil
		},
	}
}
