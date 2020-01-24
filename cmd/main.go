package main

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/login"
	"github.com/ZupIT/ritchie-cli/pkg/tree"
	ruser "github.com/ZupIT/ritchie-cli/pkg/user"
	"github.com/ZupIT/ritchie-cli/pkg/workspace"
)

const (
	ritchieHomePattern = "%s/.rit"
)

func main() {
	if env.Prod != env.Environment {
		log.Printf("Running Ritchie using %v mode. Url: %v\n\n", env.Environment, env.ServerUrl)
	}

	// get user home
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	ritchieHomePath := fmt.Sprintf(ritchieHomePattern, usr.HomeDir)

	//deps
	gitRepoManager := git.NewRepoManager()
	loginManager := login.NewDefaultManager(ritchieHomePath, env.ServerUrl, &http.Client{})
	treeManager := tree.NewDefaultManager(ritchieHomePath, env.ServerUrl, &http.Client{}, loginManager)
	credManager := credential.NewDefaultManager(env.ServerUrl, &http.Client{}, loginManager)
	userManager := ruser.NewDefaultManager(env.ServerUrl, &http.Client{}, loginManager)
	workspaceManager := workspace.NewDefaultManager(ritchieHomePath, env.ServerUrl, &http.Client{}, treeManager, gitRepoManager, credManager, loginManager)
	autocomplete := autocomplete.NewDefaultManager(env.ServerUrl, http.DefaultClient, loginManager)

	credResolver := envcredential.NewResolver(credManager)
	envResolvers := make(map[string]env.Resolver)
	envResolvers[env.Credential] = credResolver
	formulaManager := formula.NewDefaultManager(ritchieHomePath, envResolvers)

	//cmd tree
	treeBuilder := cmd.NewTreeBuilder(treeManager, workspaceManager, credManager, formulaManager, loginManager, userManager, autocomplete)
	rootCmd, err := treeBuilder.BuildTree()
	if err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
}
