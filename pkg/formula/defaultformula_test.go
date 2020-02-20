package formula

import (
	"fmt"
	"net/http"
	"os/user"
	"testing"
)

const (
	ritchieHomePattern = "%s/.rit"
)

func TestManagerMock_Run(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	ritchieHomePath := fmt.Sprintf(ritchieHomePattern, usr.HomeDir)

	f := NewDefaultManager(ritchieHomePath, nil, http.DefaultClient)

	definition := Definition{
		Path:    "kafka",
		Bin:     "kafka-${so}",
		Config:  "list-topic-config.json",
		RepoUrl: "localhost:9090/ritchie/formulas/kafka",
	}

	err = f.Run(definition)

	fmt.Println(err)
	// /formulas/github/zup-webhook/bin/zup-webhook-linux.zip
}
