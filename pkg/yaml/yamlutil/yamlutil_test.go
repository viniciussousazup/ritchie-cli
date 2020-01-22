package yamlutil

import (
	"testing"

	"github.com/matryer/is"
)

func TestReadValues(t *testing.T) {
	in := []byte(`
# Formula for k8s
formula: "Kubernetes"
version: "1.0"
dependencies:
  - "python:2.7"
  - "bash:4.0"
  - "terraform:0.12"
  - "ansible:2.0" 
entries:
  size : "big"
  auto_update: "false"
two:
  layer:
    statement: "statement"
    test: "test"
`)

	out := Values{}
	out["formula"] = "Kubernetes"
	out["version"] = "1.0"
	out["dependencies"] = []interface{}{"python:2.7", "bash:4.0", "terraform:0.12", "ansible:2.0"}
	out["entries"] = map[string]interface{}{
		"size":        "big",
		"auto_update": "false",
	}
	out["two"] = map[string]interface{}{
		"layer": map[string]interface{}{
			"statement": "statement",
			"test":      "test",
		},
	}

	is := is.New(t)

	actual, err := ReadValues(in)
	is.NoErr(err)
	is.Equal(out, actual)
}
