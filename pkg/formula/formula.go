package formula

import (
	"fmt"
	"strings"
)

const (
	// GitURL git url
	GitURL = "https://github.com/ZupIT/ritchie-formulas.git"

	// PathPattern path of formulas on local machine
	PathPattern = "%s/formulas/%s"

	// Dir Formulas
	DirFormula = "/formulas"

	// DefaultConfig is a default config file named 'config.json' for formulas
	DefaultConfig = "config.json"

	// ConfigPattern path of formula config file on local machine
	ConfigPattern = "%s/%s"

	// CommandEnv identify a COMMAND environment variable for new formula pattern.
	// This command was read in config.json file
	CommandEnv = "COMMAND"

	// BinPattern path of formula bin file on local machine
	BinPattern = "%s/bin/%s%s"
	windows    = "windows"

	// EnvPattern pattern to build envs
	EnvPattern = "%s=%s"

	// Cache pattern
	CachePattern = "%s/.%s.cache"

	DefaultCacheNewLabel = "Type new value?"

	DefaultCacheQtd = 5
)

// Config type that represents formula config
type Config struct {
	Command     string  `json:"command"`
	Description string  `json:"description"`
	Inputs      []Input `json:"inputs"`
}

// Input type that represents input config
type Input struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Default string   `json:"default"`
	Label   string   `json:"label"`
	Items   []string `json:"items"`
	Cache   Cache    `json:"cache"`
}

type Cache struct {
	Active   bool   `json:"active"`
	Qtd      int    `json:"qtd"`
	NewLabel string `json:"newLabel"`
}

// Definition type that represents a Formula
type Definition struct {
	Path    string
	Bin     string
	Config  string
	RepoUrl string
}

// FormulaPath builds the formula path from ritchie home
func (d *Definition) FormulaPath(home string) string {
	return fmt.Sprintf(PathPattern, home, d.Path)
}

// BinPath builds the bin path from formula path
func (d *Definition) BinPath(formula string, so string) string {
	suffix := ""
	if so == windows {
		suffix = ".exe"
	}
	binSO := strings.ReplaceAll(d.Bin, "${so}", so)
	return fmt.Sprintf(BinPattern, formula, binSO, suffix)
}

// ConfigPath builds the config path from formula path
func (d *Definition) ConfigPath(formula string) string {
	if d.Config != "" {
		return fmt.Sprintf(ConfigPattern, formula, d.Config)
	}
	return fmt.Sprintf(ConfigPattern, formula, DefaultConfig)
}

//go:generate $GOPATH/bin/moq -out mock_formulamanager.go . Manager

// Manager is an interface that we can use to perform formula operations
type Manager interface {
	Run(def Definition) error
}
