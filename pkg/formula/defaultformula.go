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
	fPath := def.FormulaPath(d.ritchieHome)

	var config *Config
	cName := def.ConfigName()
	cPath := def.ConfigPath(fPath, cName)
	if !fileutil.Exists(cPath) {
		if err := d.downloadConfig(def.ConfigUrl(cName), fPath, cName); err != nil {
			return err
		}
	}
	configFile, err := ioutil.ReadFile(cPath)
	if err != nil {
		return err
	}
	config = &Config{}
	if err := json.Unmarshal(configFile, config); err != nil {
		return err
	}

	bName := def.BinName()
	bPath := def.BinPath(fPath)
	bFilePath := def.BinFilePath(bPath, bName)
	if !fileutil.Exists(bFilePath) {
		zipFile, err := d.downloadFormulaBin(def.BinUrl(), bPath, bName)
		if err != nil {
			return err
		}

		if err := d.unzipFile(zipFile, bPath); err != nil {
			return err
		}
	}

	cmd := exec.Command(bFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := d.inputs(cmd, fPath, config); err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

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
				inputval, err = d.loadInputValList(items, input)
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
			e := fmt.Sprintf(EnvPattern, strings.ToUpper(input.Name), inputval)
			if i == 0 {
				cmd.Env = append(os.Environ(), e)
			} else {
				cmd.Env = append(cmd.Env, e)
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

func (d *defaultManager) loadInputValList(items []string, input Input) (string, error) {
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
	log.Println("Starting download formula...")

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("the formula bin not found")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return "", errors.New("the formula bin not found")
	default:
		return "", errors.New("unknown error when downloading your formula")
	}

	file := fmt.Sprintf("%s/%s.zip", destPath, binName)

	if err := fileutil.CreateIfNotExists(destPath, 0755); err != nil {
		return "", err
	}
	out, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	log.Println("Download formula done.")
	return file, nil
}

func (d *defaultManager) downloadConfig(url, destPath, configName string) error {
	log.Println("Starting download config file...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return errors.New("the config file not found")
	default:
		return errors.New("unknown error when downloading your config file")
	}

	file := fmt.Sprintf("%s/%s", destPath, configName)

	if err := fileutil.CreateIfNotExists(destPath, 0755); err != nil {
		return err
	}

	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	log.Println("Download config file done.")
	return nil
}

func (d *defaultManager) unzipFile(filename, destPath string) error {
	log.Println("Installing the formula...")

	if err := fileutil.CreateIfNotExists(destPath, 0655); err != nil {
		return err
	}
	if err := fileutil.Unzip(filename, destPath); err != nil {
		return err
	}
	if err := fileutil.RemoveFile(filename); err != nil {
		return err
	}

	log.Println("Formula installation done.")
	return nil
}
