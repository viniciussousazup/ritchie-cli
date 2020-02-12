package metrics

// Definition type that represents a metric use
type CmdUse struct {
	Username     string `json:"username"`
	Cmd	     string `json:"command"`
}

//go:generate $GOPATH/bin/moq -out mock_metricsmanager.go . Manager

// Manager is an interface that we can use to perform user operations
type Manager interface {
	SendCommand(cmdUse *CmdUse) error
}
