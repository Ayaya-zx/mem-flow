package inmem

import (
	"fmt"
	"sync"

	"github.com/Ayaya-zx/mem-flow/internal/common"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
)

// InmemUserRepository is an in-memory implementation of users repository.
// It is safe for concurent use by multiple goroutines.
type InmemUserRepository struct {
	m     sync.Mutex
	users map[string]*entity.User
}

func NewInmemUserRepository() *InmemUserRepository {
	return &InmemUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *InmemUserRepository) AddUser(u *entity.User) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.users[u.Name]; ok {
		return common.UserAlreadyExistsError(
			fmt.Sprintf("user with name %s already exists",
				u.Name),
		)
	}
	r.users[u.Name] = u
	return nil
}

func (r *InmemUserRepository) GetUser(name string) (*entity.User, error) {
	r.m.Lock()
	defer r.m.Unlock()
	u, ok := r.users[name]
	if !ok {
		return nil, common.UserNotExistError(
			fmt.Sprintf(
				"user with name %s does not exist", name,
			),
		)
	}
	return u, nil
}

func (r *InmemUserRepository) RemoveUser(name string) (err error) {
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.users, name)
	return nil
}
