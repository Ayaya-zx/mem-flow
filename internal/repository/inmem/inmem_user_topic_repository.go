package inmem

import (
	"fmt"
	"sync"

	repo "github.com/Ayaya-zx/mem-flow/internal/repository"
)

// InmemUserTopicRepository is an in-memory implementation
// of user topics repository. It is safe for concurent use
// by multiple goroutines.
type InmemUserTopicRepository struct {
	m             sync.Mutex
	userTopicRepo map[string]*repo.TopicRepository
}

func NewInmemUserTopicRepository() *InmemUserTopicRepository {
	return &InmemUserTopicRepository{
		userTopicRepo: make(map[string]*repo.TopicRepository),
	}
}

func (r *InmemUserTopicRepository) GetUserTopicRepository(name string) (*repo.TopicRepository, error) {
	r.m.Lock()
	defer r.m.Unlock()
	topicRepo, ok := r.userTopicRepo[name]
	if !ok {
		return nil, repo.UserTopicRepositoryNotExistError(
			fmt.Sprintf("topic repository for user %s does not exist", name),
		)
	}
	return topicRepo, nil
}

func (r *InmemUserTopicRepository) AddUserTopicRepository(name string, topicRepo *repo.TopicRepository) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.userTopicRepo[name]; ok {
		return repo.UserTopicRepositoryAlreadyExistsError(
			fmt.Sprintf("topic repository for user %s already exists", name),
		)
	}
	r.userTopicRepo[name] = topicRepo
	return nil
}
