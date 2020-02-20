package formula

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	client       *http.Client
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(ritchieHome string, ee env.Resolvers, c *http.Client) *defaultManager {
	return &defaultManager{ritchieHome: ritchieHome, envResolvers: ee, client: c}
}

// Run default implementation of function Manager.Run
func (d *defaultManager) Run(def Definition) error {
	formulaPath := def.FormulaPath(d.ritchieHome)

	var config *Config
	configName := def.ConfigName()
	configPath := def.ConfigPath(formulaPath, configName)
	if !fileutil.Exists(configPath) {
		if err := d.downloadConfig(def.ConfigUrl(configName), formulaPath, configName); err != nil {
			return err
		}
	}
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	config = &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return err
	}

	binName := def.BinName()
	binPath := def.BinPath(formulaPath)
	binFilePath := def.BinFilePath(binPath, binName)
	if !fileutil.Exists(binFilePath) {
		zipFile, err := d.downloadFormulaBin(def.BinUrl(), binPath, binName)
		if err != nil {
			return err
		}

		err = d.unzipFile(zipFile, binPath)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(binFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = d.inputs(cmd, formulaPath, config)
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}

func (d *defaultManager) inputs(cmd *exec.Cmd, formulaPath string, config *Config) error {
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
			return nil
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


func (d *defaultManager) downloadFormulaBin(url, destPath, binName string) (string, error) {
	log.Println("Starting download zip file.")

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unknown error")
	}

	file := fmt.Sprintf("%s/%s.zip", destPath, binName)

	err = fileutil.CreateIfNotExists(destPath, 0755)
	if err != nil {
		return "", err
	}
	out, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	log.Println("Download zip file done.")
	return file, nil
}

func (d *defaultManager) downloadConfig(url, destPath, configName string) error {
	log.Println("Starting download config file.")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unknown error")
	}

	file := fmt.Sprintf("%s/%s", destPath, configName)

	err = fileutil.CreateIfNotExists(destPath, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	log.Println("Download zip file done.")
	return nil
}

>>>>>>> Stashed changes
func (d *defaultManager) unzipFile(filename, destPath string) error {
	log.Println("Unzip files S3...")

	_ = fileutil.CreateIfNotExists(destPath, 0655)
	err := fileutil.Unzip(filename, destPath)
	if err != nil {
		return err
	}
	err = fileutil.RemoveFile(filename)
	if err != nil {
		return err
	}
	log.Println("Unzip S3 done.")
	return nil
}
