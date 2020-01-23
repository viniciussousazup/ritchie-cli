package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/user"
	"github.com/spf13/cobra"
	"log"
)

type deleteUserCmd struct {
	userManager user.Manager
}

func NewDeleteUserCmd(userManager user.Manager) *cobra.Command {
	o := &deleteUserCmd{userManager}

	return &cobra.Command{
		Use:   "user",
		Short: "Delete user",
		Long:  `Delete user of the organization`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt()
		},
	}
}

func (o *deleteUserCmd) prompt() error {

	un, err := prompt.String("Username: ", true)
	if err != nil {
		return err
	}
	e, err := prompt.String("Email: ", true)
	if err != nil {
		return err
	}

	u := &user.Definition{
		Email:    e,
		Username: un,
	}

	if d, err := prompt.ListBool("Are you sure want to delete this user?", []string{"yes", "no"}); err != nil {
		return err
	} else if !d {
		return nil
	}

	fmt.Println("Deleting user...")
	err = o.userManager.Delete(u)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("User %s deleted!", u.Username))

	return nil
}
