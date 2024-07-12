package inmem

import (
	"testing"
)

func TestAddAndGet(t *testing.T) {
	store := NewInmemTopicStore()
	nextId := store.nextId

	// Add task
	id, err := store.AddTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}
	if len(store.topics) != 1 {
		t.Fatalf("got len(store.topics) = %d; want 1", len(store.topics))
	}
	if _, ok := store.topicTitles["MyTopic"]; !ok {
		t.Fatalf("got store.topicTitles[\"MyTopic\"] = false; want true")
	}
	if store.nextId == nextId {
		t.Fatalf("got store.nextId = %d; want %d", store.nextId, store.nextId+1)
	}

	// Get added task
	topic, err := store.GetTopic(id)
	if err != nil {
		t.Fatal(err)
	}
	if topic.Id != id {
		t.Errorf("got topic.Id = %d; id = %d", topic.Id, id)
	}
	if topic.Title != "MyTopic" {
		t.Errorf("got topic.Title = %s; want \"MyTopic\"", topic.Title)
	}
}

func TestAddEmpty(t *testing.T) {
	store := NewInmemTopicStore()
	id := store.nextId

	// Adding a task with an empty title is prohibited
	_, err := store.AddTopic("")
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if len(store.topics) != 0 {
		t.Errorf("got len(store.topics) = %d; want 0", len(store.topics))
	}
	if store.nextId != id {
		t.Errorf("got store.nextId = %d, want %d", store.nextId, id)
	}
}

func TestAddSameTitleTwice(t *testing.T) {
	store := NewInmemTopicStore()

	store.AddTopic("MyTopic")
	nextId := store.nextId

	// We should not be able to add a task with same name twice
	_, err := store.AddTopic("MyTopic")
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if store.nextId != nextId {
		t.Errorf("got store.nextId = %d, want %d", store.nextId, nextId)
	}
}

func TestRemove(t *testing.T) {
	store := NewInmemTopicStore()

	id, err := store.AddTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}
	store.RemoveTopic(id)
	if len(store.topics) != 0 {
		t.Errorf("got len(store.topics) = %d; want 0", len(store.topics))
	}
	if len(store.topicTitles) != 0 {
		t.Errorf("got len(store.topicTitles) = %d; want 0", len(store.topicTitles))
	}
}

func TestRepeated(t *testing.T) {
	store := NewInmemTopicStore()

	id, err := store.AddTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}

	topic, err := store.GetTopic(id)
	if err != nil {
		t.Fatal(err)
	}

	nextRep := topic.NextRepeat
	err = store.TopicRepeated(id)
	if err != nil {
		t.Fatal(err)
	}
	if nextRep == topic.NextRepeat {
		t.Errorf("want topic.NextRepeat change; got %v", topic.NextRepeat)
	}

	err = store.TopicRepeated(100)
	if err == nil {
		t.Errorf("got nil; want err")
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

	store := NewInmemTopicStore()
	for _, test := range tests {
		store.AddTopic(test.title)
	}

	for _, test := range tests {
		topic, err := store.GetTopic(test.wantId)
		if err != nil {
			t.Error(err)
			continue
		}
		if topic.Title != test.title {
			t.Errorf("on id = %d want title %s; got %s",
				test.wantId, test.title, topic.Title)
		}
	}

	_, err := store.GetTopic(100)
	if err == nil {
		t.Errorf("got nil; want err")
	}
}
