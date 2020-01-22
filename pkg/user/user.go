package user

// Definition type that represents a user of the organization
type Definition struct {
	Organization string `json:"organization"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

//go:generate $GOPATH/bin/moq -out mock_usermanager.go . Manager

// Manager is an interface that we can use to perform user operations
type Manager interface {
	Create(user *Definition) error
	Delete(user *Definition) error
}
