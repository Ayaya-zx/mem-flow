package repository

import "github.com/Ayaya-zx/mem-flow/internal/entity"

// UserRepository is a representation of users repository.
type UserRepository interface {
	// AddUser adds a user to the repository.
	AddUser(u *entity.User) error
	// GetUser return user by name.
	GetUser(name string) (*entity.User, error)
	// RemoveUser deletes user from the repository by id.
	RemoveUser(name string) error
}
