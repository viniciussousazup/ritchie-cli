package sliceutil

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

var (
	coreCmds []string
)

func TestMain(m *testing.M) {
	coreCmds = []string{"root", "init", "set", "version", "credential"}
	os.Exit(m.Run())
}

func TestContains(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		in  string
		out bool
	}{
		{"init", true},
		{"notfound", false},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			contains := Contains(coreCmds, test.in)
			is.Equal(contains, test.out)
		})
	}
}
