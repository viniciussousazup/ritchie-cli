package formula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// defaultManager is a default implementation of Manager interface
type defaultManager struct {
	ritchieHome  string
	envResolvers env.Resolvers
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(ritchieHome string, ee env.Resolvers) *defaultManager {
	return &defaultManager{ritchieHome: ritchieHome, envResolvers: ee}
}

// Run default implementation of function Manager.Run
func (d *defaultManager) Run(def Definition) error {
	formulaPath := def.FormulaPath(d.ritchieHome)

	var config *Config
	configFile := def.ConfigPath(formulaPath)
	if fileutil.Exists(configFile) {
		configFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		config = &Config{}
		err = json.Unmarshal(configFile, config)
		if err != nil {
			return err
		}
	}

	so := runtime.GOOS
	bin := def.BinPath(formulaPath, so)
	cmd := exec.Command(bin)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if config != nil {
		for i, input := range config.Inputs {
			var err error
			var inputval string
			var valbool bool
			items, err := d.loadItems(input, formulaPath)
			if err != nil {
				return err
			}
			switch itype := input.Type; itype {
			case "text":
				if items != nil {
					inputval, err = d.loadInputValList(items, input, formulaPath)
				} else {
					validate := input.Default == ""
					inputval, err = prompt.String(input.Label, validate)
					if inputval == "" {
						inputval = input.Default
					}
				}
			case "bool":
				valbool, err = prompt.ListBool(input.Label, items)
				inputval = strconv.FormatBool(valbool)
			default:
				inputval, err = d.resolveIfReserved(input)
				if err != nil {
					log.Fatalf("Fail to resolve input: %v, verify your credentials. [try using set credential]", input.Type)
				}
			}

			if err != nil {
				return err
			}

			if inputval != "" {
				d.persistCache(formulaPath, inputval, input, items)
				env := fmt.Sprintf(EnvPattern, strings.ToUpper(input.Name), inputval)
				if i == 0 {
					cmd.Env = append(os.Environ(), env)
				} else {
					cmd.Env = append(cmd.Env, env)
				}
			}
		}
		if config.Command != "" {
			command := fmt.Sprintf(EnvPattern, CommandEnv, config.Command)
			cmd.Env = append(cmd.Env, command)
		}
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait()
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}

func (d *defaultManager) persistCache(formulaPath, inputVal string, input Input, items []string) {
	cachePath := fmt.Sprintf(CachePattern, formulaPath, strings.ToUpper(input.Name))
	if input.Cache.Active {
		if items == nil {
			items = []string{inputVal}
		} else {
			for i, item := range items {
				if item == inputVal { // Delete input to list
					items = append(items[:i], items[i+1:]...)
					break
				}
			}
			items = append([]string{inputVal}, items...)
		}
		qtd := DefaultCacheQtd
		if input.Cache.Qtd != 0 {
			qtd = input.Cache.Qtd
		}
		if len(items) > qtd {
			items = items[0:qtd]
		}
		itemsBytes, _ := json.Marshal(items)
		fileutil.WriteFile(cachePath, itemsBytes)
	}
}

func (d *defaultManager) loadInputValList(items []string, input Input, formulaPath string) (string, error) {
	newLabel := DefaultCacheNewLabel
	if input.Cache.Active {
		if input.Cache.NewLabel != "" {
			newLabel = input.Cache.NewLabel
		}
		items = append(items, newLabel)
	}
	inputval, err := prompt.List(input.Label, items)
	if inputval == newLabel {
		validate := input.Default == ""
		inputval, err = prompt.String(input.Label, validate)
		if inputval == "" {
			inputval = input.Default
		}
	}
	return inputval, err
}

func (d *defaultManager) loadItems(input Input, formulaPath string) ([]string, error) {
	if input.Cache.Active {
		cachePath := fmt.Sprintf(CachePattern, formulaPath, strings.ToUpper(input.Name))
		if fileutil.Exists(cachePath) {
			fileBytes, err := fileutil.ReadFile(cachePath)
			if err != nil {
				return nil, err
			}
			var items []string
			err = json.Unmarshal(fileBytes, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		} else {
			itemsBytes, err := json.Marshal(input.Items)
			if err != nil {
				return nil, err
			}
			err = fileutil.WriteFile(cachePath, itemsBytes)
			if err != nil {
				return nil, err
			}
			return input.Items, nil
		}
	} else {
		return input.Items, nil
	}
}

func (d *defaultManager) resolveIfReserved(input Input) (string, error) {
	s := strings.Split(input.Type, "_")
	resolver := d.envResolvers[s[0]]
	if resolver != nil {
		return resolver.Resolve(input.Type)
	}
	return "", nil
}
