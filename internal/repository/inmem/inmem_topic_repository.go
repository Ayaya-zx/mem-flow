package inmem

import (
	"fmt"
	"sync"
	"time"

	"github.com/Ayaya-zx/mem-flow/internal/entity"
	repo "github.com/Ayaya-zx/mem-flow/internal/repository"
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
		nextId:      1,
		topics:      make(map[int]*entity.Topic),
		topicTitles: make(map[string]struct{}),
	}
}

func (ts *InmemTopicRepository) AddTopic(title string) (int, error) {
	if title == "" {
		return 0, repo.TopicTitleError("topic's title is empty")
	}
	if _, ok := ts.topicTitles[title]; ok {
		return 0, repo.TopicTitleError(fmt.Sprintf(
			"topic with title %s already exists",
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
	topic.Id = ts.nextId
	ts.topics[ts.nextId] = topic
	ts.nextId++
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

func (ts *InmemTopicRepository) GetTopic(id int) (*entity.Topic, error) {
	ts.m.Lock()
	defer ts.m.Unlock()
	t, ok := ts.topics[id]
	if !ok {
		return nil, repo.TopicNotExistsError(
			fmt.Sprintf("topic with id %d does not exist", id))
	}
	return t, nil
}

func (ts *InmemTopicRepository) TopicRepeated(id int) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	t, ok := ts.topics[id]
	if !ok {
		return repo.TopicNotExistsError(
			fmt.Sprintf("topic with id %d does not exist", id))
	}
	t.Repeat()
	return nil
}
