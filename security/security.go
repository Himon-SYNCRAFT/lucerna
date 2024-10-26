package security

type User interface {
	GetUserIdentifier() string
	GetRoles() []string
}

type Security interface {
	GetUser() (*User, error)
}
