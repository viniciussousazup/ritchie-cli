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
	// ConfigPattern path of formula config file on local machine
	ConfigPattern = "%s/config.json"
	// BinPattern path of formula bin file on local machine
	BinPattern = "%s/bin/%s%s"
	windows    = "windows"
	// EnvPattern pattern to build envs
	EnvPattern = "%s=%s"
)

// Config type that represents formula config
type Config struct {
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
}

// Definition type that represents a Formula
type Definition struct {
	Path string
	Bin  string
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
	return fmt.Sprintf(ConfigPattern, formula)
}

//go:generate $GOPATH/bin/moq -out mock_formulamanager.go . Manager

// Manager is an interface that we can use to perform formula operations
type Manager interface {
	Run(def Definition) error
}