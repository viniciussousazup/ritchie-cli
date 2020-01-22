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

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// DefaultManager is a default implementation of Manager interface
type DefaultManager struct {
	ritchieHome  string
	envResolvers map[string]env.Resolver
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(ritchieHome string, envResolvers map[string]env.Resolver) *DefaultManager {
	return &DefaultManager{ritchieHome, envResolvers}
}

//Run default implementation of function Manager.Run
func (d *DefaultManager) Run(def Definition) error {
	formulaPath := def.FormulaPath(d.ritchieHome)

	var config *Config
	configFile := def.ConfigPath(formulaPath)
	if fileutil.Exists(configFile) {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		config = &Config{}
		err = json.Unmarshal([]byte(file), config)
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
			items := input.Items

			switch itype := input.Type; itype {
			case "text":
				if items != nil {
					inputval, err = prompt.List(input.Label, items)
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
				env := fmt.Sprintf(EnvPattern, strings.ToUpper(input.Name), inputval)
				if i == 0 {
					cmd.Env = append(os.Environ(), env)
				} else {
					cmd.Env = append(cmd.Env, env)
				}
			}
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

func (d *DefaultManager) resolveIfReserved(input Input) (string, error) {
	s := strings.Split(input.Type, "_")
	resolver := d.envResolvers[s[0]]
	if resolver != nil {
		return resolver.Resolve(input.Type)
	}
	return "", nil
}
