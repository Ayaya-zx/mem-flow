package repository

// UserTopicRepository stores TopicRepository instances associated
// with user names.
type UserTopicRepository interface {
	// GetUserTopicRepository returns TopicRepository instance associated
	// with the user with the given name.
	GetUserTopicRepository(name string) (TopicRepository, error)
}
