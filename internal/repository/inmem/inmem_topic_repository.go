package inmem

import (
	"fmt"
	"sync"
	"time"

	"github.com/Ayaya-zx/mem-flow/internal/common"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
)

// InmemTopicRepository is an in-memory implementation of topics repository.
// It is safe for concurent use by multiple goroutines.
type InmemTopicRepository struct {
	m      sync.Mutex
	topics map[string]*entity.Topic
}

func NewInmemTopicRepository() *InmemTopicRepository {
	return &InmemTopicRepository{
		topics: make(map[string]*entity.Topic),
	}
}

func (ts *InmemTopicRepository) AddTopic(title string) error {
	if title == "" {
		return common.TopicTitleError("topic's title is empty")
	}
	if _, ok := ts.topics[title]; ok {
		return common.TopicTitleError(fmt.Sprintf(
			"topic %s already exists",
			title,
		))
	}

	topic := new(entity.Topic)
	topic.Title = title
	topic.Created = time.Now()
	topic.LastRepeated = topic.Created
	topic.NextRepeat = time.Now().Add(time.Duration(20 * time.Minute))

	ts.m.Lock()
	defer ts.m.Unlock()
	ts.topics[title] = topic
	return nil
}

func (ts *InmemTopicRepository) RemoveTopic(title string) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	delete(ts.topics, title)
	return nil
}

func (ts *InmemTopicRepository) GetAllTopics() ([]*entity.Topic, error) {
	ts.m.Lock()
	defer ts.m.Unlock()
	res := make([]*entity.Topic, 0, len(ts.topics))
	for _, t := range ts.topics {
		res = append(res, t)
	}
	return res, nil
}

func (ts *InmemTopicRepository) GetTopic(title string) (*entity.Topic, error) {
	ts.m.Lock()
	defer ts.m.Unlock()
	t, ok := ts.topics[title]
	if !ok {
		return nil, common.TopicNotExistsError(
			fmt.Sprintf("topic %s does not exist", title))
	}
	return t, nil
}
