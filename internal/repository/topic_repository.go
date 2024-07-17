package repository

import "github.com/Ayaya-zx/mem-flow/internal/entity"

type (
	// TopicTitleError should be returned by the method AddTopic
	// of the type TopicStore if a topic with a given title cannot
	// be added to the repository.
	TopicTitleError string
	// TopicNotExistsError is returned by the GetTopic method of type
	// TopicStore, if a topic with the given id is not in the repository.
	TopicNotExistsError string
)

func (e TopicTitleError) Error() string {
	return string(e)
}

func (e TopicNotExistsError) Error() string {
	return string(e)
}

type TopicRepository interface {
	// AddTopic adds a topic with a given title to the repository.
	AddTopic(title string) (int, error)
	// RemoveTopic deletes a topic from the repository by id.
	RemoveTopic(id int) error
	// GetAllTopics returns all topics stored at the repository.
	GetAllTopics() ([]*entity.Topic, error)
	// GetTopic returns topic by id.
	GetTopic(id int) (*entity.Topic, error)
}
