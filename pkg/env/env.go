package env

const (
	// Credential resolver
	Credential  = "CREDENTIAL"
	Prod		= "prod"
)

var (
	//Url Server
	ServerUrl = "https://ritchie-server.zup.io"
	// Environment
	Environment = Prod
)

//go:generate $GOPATH/bin/moq -out mock_envresolver.go . Resolver

type Resolvers map[string]Resolver

// Resolver is an interface that we can use to resolve reserved envs
type Resolver interface {
	Resolve(name string) (string, error)
}
