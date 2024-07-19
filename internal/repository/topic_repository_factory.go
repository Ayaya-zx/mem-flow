package repository

type TopicRepositoryFactory interface {
	CreateTopicRepository() (TopicRepository, error)
}
