package env

const (
	// Credential resolver
	Credential  = "CREDENTIAL"
	Dev			= "dev"
	Prod		= "prod"
)

var (
	//Url Server
	ServerUrl = "https://ritchie-server.itiaws.dev"
	// Environment
	Environment = Prod
)

//go:generate $GOPATH/bin/moq -out mock_envresolver.go . Resolver

type Resolvers map[string]Resolver

// Resolver is an interface that we can use to resolve reserved envs
type Resolver interface {
	Resolve(name string) (string, error)
}
