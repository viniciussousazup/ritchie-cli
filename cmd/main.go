package main

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/login"
	"github.com/ZupIT/ritchie-cli/pkg/metrics"
	"github.com/ZupIT/ritchie-cli/pkg/tree"
	ruser "github.com/ZupIT/ritchie-cli/pkg/user"
	"github.com/ZupIT/ritchie-cli/pkg/workspace"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"
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
	gitManager := git.NewDefaultManager()
	loginManager := login.NewDefaultManager(ritchieHomePath, env.ServerUrl, http.DefaultClient)
	treeManager := tree.NewDefaultManager(ritchieHomePath, env.ServerUrl, http.DefaultClient, loginManager)
	credManager := credential.NewDefaultManager(env.ServerUrl, http.DefaultClient, loginManager)
	userManager := ruser.NewDefaultManager(env.ServerUrl, http.DefaultClient, loginManager)
	workspaceManager := workspace.NewDefaultManager(ritchieHomePath, env.ServerUrl, http.DefaultClient, treeManager, gitManager, credManager, loginManager)
	autocompleteManager := autocomplete.NewDefaultManager(env.ServerUrl, http.DefaultClient, loginManager)
	metricsManager := metrics.NewDefaultManager(env.ServerUrl,&http.Client{Timeout: 2 *time.Second},loginManager)

	credResolver := envcredential.NewResolver(credManager)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver
	formulaManager := formula.NewDefaultManager(ritchieHomePath, envResolvers)

	//cmd tree
	treeBuilder := cmd.NewTreeBuilder(treeManager, workspaceManager, credManager, formulaManager, loginManager, userManager, autocompleteManager)
	rootCmd, err := treeBuilder.BuildTree()
	if err != nil {
		panic(err)
	}

	go metricsManager.SendCommand()

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}

}
