package inmem

import (
	"testing"
)

func TestAddAndGet(t *testing.T) {
	repo := NewInmemTopicRepository()
	nextId := repo.nextId

	// Add task
	id, err := repo.AddTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.topics) != 1 {
		t.Fatalf("got len(repo.topics) = %d; want 1", len(repo.topics))
	}
	if _, ok := repo.topicTitles["MyTopic"]; !ok {
		t.Fatalf("got repo.topicTitles[\"MyTopic\"] = false; want true")
	}
	if repo.nextId == nextId {
		t.Fatalf("got repo.nextId = %d; want %d", repo.nextId, repo.nextId+1)
	}

	// Get added task
	topic, err := repo.GetTopic(id)
	if err != nil {
		t.Fatal(err)
	}
	// if topic.Id != id {
	// 	t.Errorf("got topic.Id = %d; id = %d", topic.Id, id)
	// }
	if topic.Title != "MyTopic" {
		t.Errorf("got topic.Title = %s; want \"MyTopic\"", topic.Title)
	}
}

func TestAddEmpty(t *testing.T) {
	repo := NewInmemTopicRepository()
	id := repo.nextId

	// Adding a task with an empty title is prohibited
	_, err := repo.AddTopic("")
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if len(repo.topics) != 0 {
		t.Errorf("got len(repo.topics) = %d; want 0", len(repo.topics))
	}
	if repo.nextId != id {
		t.Errorf("got repo.nextId = %d, want %d", repo.nextId, id)
	}
}

func TestAddSameTitleTwice(t *testing.T) {
	repo := NewInmemTopicRepository()

	repo.AddTopic("MyTopic")
	nextId := repo.nextId

	// We should not be able to add a task with same name twice
	_, err := repo.AddTopic("MyTopic")
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if repo.nextId != nextId {
		t.Errorf("got repo.nextId = %d, want %d", repo.nextId, nextId)
	}
}

func TestRemove(t *testing.T) {
	repo := NewInmemTopicRepository()

	id, err := repo.AddTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}
	repo.RemoveTopic(id)
	if len(repo.topics) != 0 {
		t.Errorf("got len(repo.topics) = %d; want 0", len(repo.topics))
	}
	if len(repo.topicTitles) != 0 {
		t.Errorf("got len(repo.topicTitles) = %d; want 0", len(repo.topicTitles))
	}
}

func TestGet(t *testing.T) {
	var tests = []struct {
		title  string
		wantId int
	}{
		{"MyTopic1", 1},
		{"MyTopic2", 2},
		{"MyTopic3", 3},
	}

	topicRepo := NewInmemTopicRepository()
	for _, test := range tests {
		topicRepo.AddTopic(test.title)
	}

	for _, test := range tests {
		topic, err := topicRepo.GetTopic(test.wantId)
		if err != nil {
			t.Error(err)
			continue
		}
		if topic.Title != test.title {
			t.Errorf("on id = %d want title %s; got %s",
				test.wantId, test.title, topic.Title)
		}
	}

	_, err := topicRepo.GetTopic(100)
	if err == nil {
		t.Errorf("got nil; want err")
	}
}
