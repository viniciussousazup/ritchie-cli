package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"

	"github.com/spf13/cobra"
)

// SetCredentialCmd type for set credential command
type SetCredentialCmd struct {
	manager credential.Manager
}

// NewSetCredentialCmd creates a new cmd instance
func NewSetCredentialCmd(m credential.Manager) *cobra.Command {
	o := &SetCredentialCmd{m}

	return &cobra.Command{
		Use:   "credential",
		Short: "Set provider credential",
		Long:  `Set credentials for Github, Gitlab, AWS, etc.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.prompt()
		},
	}
}

func (s *SetCredentialCmd) prompt() error {
	cfg, err := s.manager.Configs()
	if err != nil {
		return err
	}
	providers := make([]string, 0, len(cfg))
	for k := range cfg {
		providers = append(providers, k)
	}

	typ, err := prompt.List("Type: ", []string{credential.Me, credential.Admin})
	if err != nil {
		return err
	}

	username := "me"
	if typ == credential.Admin {
		username, err = prompt.String("Set credential for user [username] ", true)
		if err != nil {
			return err
		}
	}

	provider, err := prompt.List("Provider: ", providers)
	if err != nil {
		return err
	}

	values := make(map[string]string)
	fields := cfg[provider]
	for _, f := range fields {
		var val string
		var err error
		fname := strings.ToLower(f.Field)
		lab := fmt.Sprintf("%s: ", strings.Title(f.Field))
		if f.Type == prompt.PasswordType {
			val, err = prompt.Password(lab)
		} else {
			val, err = prompt.String(lab, true)
		}
		if err != nil {
			return err
		}
		values[fname] = val
	}

	cred := &credential.Secret{
		Username:   username,
		Credential: values,
		Provider:   provider,
	}

	err = s.manager.Save(cred)
	if err != nil {
		return err
	}

	log.Println("Credential saved!")
	return nil
}
