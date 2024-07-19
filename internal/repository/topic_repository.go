package repository

import "github.com/Ayaya-zx/mem-flow/internal/entity"

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
