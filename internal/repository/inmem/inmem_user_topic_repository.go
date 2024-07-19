package inmem

import (
	"sync"

	repo "github.com/Ayaya-zx/mem-flow/internal/repository"
)

// InmemUserTopicRepository is an in-memory implementation
// of user topics repository. It is safe for concurent use
// by multiple goroutines.
type InmemUserTopicRepository struct {
	m                sync.Mutex
	userTopicRepo    map[string]repo.TopicRepository
	topicRepoFactory repo.TopicRepositoryFactory
}

func NewInmemUserTopicRepository(topicRepoFactory repo.TopicRepositoryFactory) *InmemUserTopicRepository {
	return &InmemUserTopicRepository{
		userTopicRepo:    make(map[string]repo.TopicRepository),
		topicRepoFactory: topicRepoFactory,
	}
}

func (r *InmemUserTopicRepository) GetUserTopicRepository(name string) (repo.TopicRepository, error) {
	r.m.Lock()
	defer r.m.Unlock()
	topicRepo, ok := r.userTopicRepo[name]
	if !ok {
		var err error
		topicRepo, err = r.topicRepoFactory.CreateTopicRepository()
		if err != nil {
			return nil, err
		}
		r.userTopicRepo[name] = topicRepo
	}
	return topicRepo, nil
}
