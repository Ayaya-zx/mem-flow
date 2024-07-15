package repository

import "github.com/Ayaya-zx/mem-flow/internal/entity"

// TopicTitleError should be returned by the method AddTopic
// of the type TopicStore if a topic with a given title cannot
// be added to the storage.
type TopicTitleError string

func (e TopicTitleError) Error() string {
	return string(e)
}

// TopicNotExistsError should be returned by method GetTopic of
// the type TopicStore if a topic with given id is not in the storage.
type TopicNotExistsError string

func (e TopicNotExistsError) Error() string {
	return string(e)
}

// TopicRepository is a representation of topics storage.
type TopicRepository interface {
	// AddTopic adds a topic with a given title to the storage.
	AddTopic(string) (int, error)
	// RemoveTopic deletes a topic from the storage by id.
	RemoveTopic(int) error
	// GetAllTopics returns all topics stored at the storage.
	GetAllTopics() ([]*entity.Topic, error)
	// GetTopic returns topic by id.
	GetTopic(int) (*entity.Topic, error)
}
