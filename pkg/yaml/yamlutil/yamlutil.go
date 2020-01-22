package yamlutil

import (
	"bytes"
	"log"
	"text/template"

	"sigs.k8s.io/yaml"
)

// Values map type that represents a yaml fields/values
type Values map[string]interface{}

// Tpl type that represents a yaml template to be compiled
type Tpl struct {
	Name    string
	Content string
	Values  interface{}
}

// YAML converts a yaml map to byte array
func (v Values) YAML() ([]byte, error) {
	b, err := yaml.Marshal(v)
	return b, err
}

// ReadValues reads a byte array to map
func ReadValues(data []byte) (vals Values, err error) {
	err = yaml.Unmarshal(data, &vals)
	if len(vals) == 0 {
		vals = Values{}
	}

	log.Printf("Values: %v\n", vals)
	return vals, err
}

// Compile compiles a template with values
func Compile(tpl Tpl) ([]byte, error) {
	var b bytes.Buffer
	tt := template.Must(template.New(tpl.Name).Parse(tpl.Content))
	err := tt.Execute(&b, tpl.Values)
	return b.Bytes(), err
}
