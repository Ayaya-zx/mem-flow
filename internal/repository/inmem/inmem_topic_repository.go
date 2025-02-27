package inmem

import (
	"fmt"
	"sync"

	"github.com/Ayaya-zx/mem-flow/internal/common"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
)

// InmemTopicRepository is an in-memory implementation of topics repository.
// It is safe for concurent use by multiple goroutines.
type InmemTopicRepository struct {
	m           sync.Mutex
	topics      map[int]*entity.Topic
	topicTitles map[string]struct{}
	nextId      int
}

func NewInmemTopicRepository() *InmemTopicRepository {
	return &InmemTopicRepository{
		topics:      make(map[int]*entity.Topic),
		topicTitles: make(map[string]struct{}),
		nextId:      1,
	}
}

func (ts *InmemTopicRepository) AddTopic(title string) (int, error) {
	if title == "" {
		return 0, common.TopicTitleError("topic's title is empty")
	}

	ts.m.Lock()
	defer ts.m.Unlock()

	if _, ok := ts.topicTitles[title]; ok {
		return 0, common.TopicTitleError(fmt.Sprintf(
			"topic %s already exists",
			title,
		))
	}

	topic := entity.NewTopic(ts.nextId, title)
	ts.nextId++
	ts.topics[topic.Id] = topic
	ts.topicTitles[title] = struct{}{}

	return topic.Id, nil
}

func (ts *InmemTopicRepository) RemoveTopic(id int) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	if topic, ok := ts.topics[id]; ok {
		delete(ts.topics, id)
		delete(ts.topicTitles, topic.Title)
	}
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

func (ts *InmemTopicRepository) GetTopicById(id int) (*entity.Topic, error) {
	ts.m.Lock()
	defer ts.m.Unlock()
	t, ok := ts.topics[id]
	if !ok {
		return nil, common.TopicNotExistsError(
			fmt.Sprintf("topic with id %d does not exist", id))
	}
	return t, nil
}
