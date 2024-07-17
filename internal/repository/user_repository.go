package repository

import "github.com/Ayaya-zx/mem-flow/internal/entity"

type (
	// UserNotExistError is returned by the GetUser method
	// of type UserRepository, if a user with the given name
	// is not in the repository.
	UserNotExistError string
	// UserAlreadyExistsError is returned by the AddUser method
	// of type UserRepository, if a user with given name already exists.
	UserAlreadyExistsError string
)

func (e UserNotExistError) Error() string {
	return string(e)
}

func (e UserAlreadyExistsError) Error() string {
	return string(e)
}

// UserRepository is a representation of users repository.
type UserRepository interface {
	// AddUser adds a user to the repository.
	AddUser(u *entity.User) error
	// GetUser return user by name.
	GetUser(name string) (*entity.User, error)
	// RemoveUser deletes user from the repository by id.
	RemoveUser(name string) error
}
