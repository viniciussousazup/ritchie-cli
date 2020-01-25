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
func NewLoginCmd(loginManager login.Manager) *cobra.Command {
	o := &loginCmd{loginManager}
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
	org, err := prompt.String("Organization: ", true)
	if err != nil {
		return err
	}
	u, err := prompt.String("Username: ", true)
	if err != nil {
		return nil
	}
	p, err := prompt.Password("Password: ")
	if err != nil {
		return nil
	}

	c := &login.Credential{
		Username:     u,
		Password:     p,
		Organization: org,
	}

	err = o.loginManager.Authenticate(c)
	if err != nil {
		return err
	}

	log.Println("Login successful!")
	return nil
}
