package repository

type (
	// UserTopicRepositoryNotExistError is returned by
	// the GetUserTopicRepository method of type UserTopicRepository,
	// if there is no topic repository associated with the user with
	// the given name.
	UserTopicRepositoryNotExistError string
	// UserTopicRepositoryAlreadyExistsError is returned by
	// the AddUserTopicRepository method of type UserTopicRepository,
	// if a topic repository associated with the user with the given
	// name already exists.
	UserTopicRepositoryAlreadyExistsError string
)

func (e UserTopicRepositoryNotExistError) Error() string {
	return string(e)
}

func (e UserTopicRepositoryAlreadyExistsError) Error() string {
	return string(e)
}

// UserTopicRepository stores TopicRepository instances associated
// with user names.
type UserTopicRepository interface {
	// GetUserTopicRepository returns TopicRepository instance associated
	// with the user with the given name.
	GetUserTopicRepository(name string) (*UserRepository, error)
	// AddUserTopicRepository associates given topic repository with
	// the given user name.
	AddUserTopicRepository(name string, topicRepository *TopicRepository) error
}
