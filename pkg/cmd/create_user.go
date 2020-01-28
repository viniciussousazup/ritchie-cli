package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/user"
)

// createUserCmd type for create user command
type createUserCmd struct {
	userManager user.Manager
}

// NewCreateUserCmd creates a new cmd instance
func NewCreateUserCmd(userManager user.Manager) *cobra.Command {
	o := &createUserCmd{userManager}

	return &cobra.Command{
		Use:   "user",
		Short: "Create user",
		Long:  `Create user of the organization`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt(args)
		},
	}
}

func (o *createUserCmd) prompt(args []string) error {
	org, err := prompt.String("Organization: ", true)
	if err != nil {
		return err
	}
	fn, err := prompt.String("First name: ", true)
	if err != nil {
		return err
	}
	ln, err := prompt.String("Last name: ", true)
	if err != nil {
		return err
	}
	e, err := prompt.Email("Email: ")
	if err != nil {
		return err
	}
	un, err := prompt.String("Username: ", true)
	if err != nil {
		return err
	}
	p, err := prompt.Password("Password: ")
	if err != nil {
		return err
	}

	u := &user.Definition{
		Organization: org,
		FirstName:    fn,
		LastName:     ln,
		Email:        e,
		Username:     un,
		Password:     p,
	}

	err = o.userManager.Create(u)
	if err != nil {
		return err
	}

	log.Println("User created!")

	return err
}
