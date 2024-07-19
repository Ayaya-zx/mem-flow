package inmem

import repo "github.com/Ayaya-zx/mem-flow/internal/repository"

type InmemTopicRepositoryFactory struct{}

func NewInmemTopicRepositoryFactory() *InmemTopicRepositoryFactory {
	return &InmemTopicRepositoryFactory{}
}

func (InmemTopicRepositoryFactory) CreateTopicRepository() (repo.TopicRepository, error) {
	return NewInmemTopicRepository(), nil
}
