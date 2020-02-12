package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/context"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/login"
	"github.com/ZupIT/ritchie-cli/pkg/tree"
	"github.com/ZupIT/ritchie-cli/pkg/user"
	"github.com/ZupIT/ritchie-cli/pkg/workspace"
)

var (
	coreCmds = []string{
		"root_init", "root_set", "root_set_credential", "root_set_context",
		"root_create", "root_create_user", "root_delete", "root_delete_user",
		"root_delete_context", "root_show", "root_show_context", "root_help",
		"root_login", "root_completion", "root_completion_zsh", "root_completion_bash"}
)

// treeBuilder type that represents the tree builder
type treeBuilder struct {
	treeManager         tree.Manager
	workspaceManager    workspace.Manager
	credManager         credential.Manager
	formulaManager      formula.Manager
	loginManager        login.Manager
	userManager         user.Manager
	autocompleteManager autocomplete.Manager
	ctxManager          context.Manager
}

// NewTreeBuilder creates new builder instance
func NewTreeBuilder(
	t tree.Manager,
	w workspace.Manager,
	c credential.Manager,
	f formula.Manager,
	l login.Manager,
	u user.Manager,
	a autocomplete.Manager,
	ct context.Manager,
) *treeBuilder {
	return &treeBuilder{
		treeManager:         t,
		workspaceManager:    w,
		credManager:         c,
		formulaManager:      f,
		loginManager:        l,
		userManager:         u,
		autocompleteManager: a,
		ctxManager:          ct,
	}
}

// BuildTree builds the tree of the commands
func (b *treeBuilder) BuildTree() (*cobra.Command, error) {
	rootCmd := NewRootCmd(b.workspaceManager)
	initCmd := NewInitCmd(b.workspaceManager)
	setCmd := NewSetCmd()
	showCmd := NewShowCmd()
	createCmd := NewCreateCmd()
	deleteCmd := NewDeleteCmd()
	loginCmd := NewLoginCmd(b.loginManager)
	setCredentialCmd := NewSetCredentialCmd(b.credManager)
	setCtxCmd := NewSetContextCmd(b.ctxManager)
	createUserCmd := NewCreateUserCmd(b.userManager)
	deleteUserCmd := NewDeleteUserCmd(b.userManager)
	deleteCtxCmd := NewDeleteContextCmd(b.ctxManager)
	autocompleteCmd := NewAutocompleteCmd()
	autocompleteZsh := NewAutocompleteZsh(b.autocompleteManager)
	autocompleteBash := NewAutocompleteBash(b.autocompleteManager)
	showCtxCmd := NewShowContextCmd(b.ctxManager)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash)
	setCmd.AddCommand(setCredentialCmd)
	setCmd.AddCommand(setCtxCmd)
	createCmd.AddCommand(createUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)
	deleteCmd.AddCommand(deleteCtxCmd)
	showCmd.AddCommand(showCtxCmd)

	rootCmd.AddCommand(initCmd, setCmd, createCmd, deleteCmd, loginCmd, autocompleteCmd, showCmd)

	treeCmd, err := b.treeManager.GetLocalTree()
	if err != nil {
		return nil, err
	} else if treeCmd != nil {
		b.processTree(treeCmd, rootCmd)
	}

	return rootCmd, nil
}

func (b *treeBuilder) processTree(treeCmd *tree.Representation, rootCmd *cobra.Command) {
	commands := make(map[string]*cobra.Command)
	commands["root"] = rootCmd

	for _, v := range treeCmd.Commands {
		cmdKey := fmt.Sprintf("%s_%s", v.Parent, v.Usage)
		if !sliceutil.Contains(coreCmds, cmdKey) {
			var annotations map[string]string

			var cmd *cobra.Command
			if v.Formula.Path != "" {
				f := v.Formula
				annotations = make(map[string]string)
				annotations["formulaPath"] = f.Path
				annotations["formulaBin"] = f.Bin
				annotations["formulaConfig"] = f.Config
				cmd = &cobra.Command{
					Use:   v.Usage,
					Short: v.Help,
					Long:  v.Help,
					RunE: func(cmd *cobra.Command, args []string) error {
						log.Printf("Running cmd %v with args %v", cmd.Use, args)
						if cmd.Annotations != nil {
							fPath := cmd.Annotations["formulaPath"]
							fBin := cmd.Annotations["formulaBin"]
							fConf := cmd.Annotations["formulaConfig"]
							frm := formula.Definition{
								Path:   fPath,
								Bin:    fBin,
								Config: fConf,
							}
							b.formulaManager.Run(frm)
						}
						return nil
					},
				}
			} else {
				cmd = &cobra.Command{
					Use:   v.Usage + " SUBCOMMAND",
					Short: v.Help,
					Long:  v.Help,
				}
			}

			if annotations != nil {
				cmd.Annotations = annotations
			}

			parentCmd := commands[v.Parent]
			parentCmd.AddCommand(cmd)
			commands[cmdKey] = cmd
		}
	}
}
