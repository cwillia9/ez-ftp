package domain

type UserRepository interface {
	FindByName(string) (User, error)
	Store(User) error
}

type User struct {
	name     string
	passhash string
}
