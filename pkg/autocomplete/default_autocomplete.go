package autocomplete

import (
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/tree"
	"net/http"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

const (
	binaryName  string = "rit"
	lineCommand string = "    commands+=(\"${command}\")"
	firstLevel  string = "root"
	bash        string = "bash"
	zsh         string = "zsh"
)

type defaultManager struct {
	serverURL   string
	httpClient  *http.Client
	ritchieHome string
}

type (
	BashCommand struct {
		LastCommand string
		RootCommand string
		Commands    string
		Level       int
	}

	Command struct {
		Content []string
		Before  string
	}
)

func NewDefaultManager(serverURL, ritchieHome string, c *http.Client) *defaultManager {
	return &defaultManager{serverURL: serverURL, httpClient: c, ritchieHome: ritchieHome}
}

func (d *defaultManager) Handle(shellName string) (string, error) {
	if !sliceutil.Contains(supportedAutocomplete(), shellName) {
		return "", errors.New("autocomplete for this terminal is not supported")
	}
	treeManager := tree.NewDefaultManager(d.ritchieHome, "", nil, nil)
	tree, err := treeManager.GetLocalTree()
	if err != nil {
		return "", err
	}
	autoCompletion := ""
	switch shellName {
	case bash:
		autoCompletion = fmt.Sprintf("%s\n%s", "#!/bin/bash", loadToBash(*tree))
	case zsh:
		autoCompletion = loadToZsh(*tree)
	}
	return autoCompletion, nil
}

func supportedAutocomplete() []string {
	return []string{bash, zsh}
}

func loadToBash(tree tree.Representation) string {
	autoComplete := autoCompletionBash
	autoComplete = strings.Replace(autoComplete, "{{BinaryName}}", binaryName, -1)
	autoComplete = strings.Replace(autoComplete, "{{DynamicCode}}", loadDynamicCommands(tree), 1)
	return autoComplete
}

func loadToZsh(tree tree.Representation) string {
	autoComplete := autoCompletionZsh
	autoComplete = strings.Replace(autoComplete, "{{BinaryName}}", binaryName, -1)
	autoComplete = strings.Replace(autoComplete, "{{AutoCompleteBash}}", loadToBash(tree), 1)
	return autoComplete
}

func loadDynamicCommands(tree tree.Representation) string {
	commands := tree.Commands
	commandString := command
	mapCommand := loadCommands(commands)
	bashCommands := loadBashCommands(mapCommand)

	allCommands := ""
	for _, bashCommand := range bashCommands {
		functionName := formatterFunctionName(bashCommand.RootCommand)
		command := strings.Replace(commandString, "{{RootCommand}}", bashCommand.RootCommand, -1)
		command = strings.Replace(command, "{{LastCommand}}", bashCommand.LastCommand, -1)
		command = strings.Replace(command, "{{FunctionName}}", functionName, -1)
		allCommands += strings.Replace(command, "{{Commands}}", bashCommand.Commands, -1)
	}
	return allCommands
}

func formatterFunctionName(functionName string) string {
	functionParts := strings.Split(functionName, "_")
	if len(functionParts) > 2 {
		functionName = functionParts[len(functionParts)-2] + "_" + functionParts[len(functionParts)-1]
	}
	return functionName
}

func loadCommands(commands []tree.Command) map[string]Command {
	commandsMap := make(map[string]Command)
	for _, command := range commands {
		addValueMap(&commandsMap, command.Parent, command.Usage, command.Parent)
	}
	commandsMapResponse := make(map[string]Command)
	for key, value := range commandsMap {
		commandsMapResponse[key] = value
		for _, v := range value.Content {
			newKey := key + "_" + v
			if _, ok := commandsMap[newKey]; !ok {
				commandsMapResponse[newKey] = Command{
					Content: nil,
					Before:  newKey,
				}
			}
		}
	}
	return commandsMapResponse
}

func loadBashCommands(mapCommands map[string]Command) []BashCommand {
	var bashCommands []BashCommand
	for key, value := range mapCommands {
		rootCommand := key
		level := len(strings.Split(key, "_"))
		commands := ""
		for _, valueEntry := range value.Content {
			commands += strings.Replace(lineCommand, "${command}", valueEntry, -1) + "\n"
		}
		if rootCommand == firstLevel {
			rootCommand = fmt.Sprintf("%s_%s", binaryName, rootCommand)
		}
		bashCommands = append(bashCommands, BashCommand{
			RootCommand: rootCommand,
			Commands:    commands,
			LastCommand: loadLastCommand(key),
			Level:       level,
		},
		)
	}
	return bashCommands
}

func loadLastCommand(rootCommand string) string {
	splitRootCommand := strings.Split(rootCommand, "_")
	return splitRootCommand[len(splitRootCommand)-1]
}

func addValueMap(mapCommand *map[string]Command, key string, value string, before string) {
	i := *mapCommand
	a := i[key]
	a.Content = append(a.Content, value)
	a.Before = before
	i[key] = a
}
