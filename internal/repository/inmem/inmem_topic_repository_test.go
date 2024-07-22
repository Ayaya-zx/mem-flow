package inmem

import (
	"testing"
)

func TestAddAndGet(t *testing.T) {
	repo := NewInmemTopicRepository()

	// Add task
	err := repo.AddTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.topics) != 1 {
		t.Fatalf("got len(repo.topics) = %d; want 1", len(repo.topics))
	}
	if _, ok := repo.topics["MyTopic"]; !ok {
		t.Fatalf("got repo.topicTitles[\"MyTopic\"] = false; want true")
	}

	// Get added task
	topic, err := repo.GetTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}
	if topic.Title != "MyTopic" {
		t.Errorf("got topic.Title = %s; want \"MyTopic\"", topic.Title)
	}
}

func TestAddEmpty(t *testing.T) {
	repo := NewInmemTopicRepository()

	// Adding a task with an empty title is prohibited
	err := repo.AddTopic("")
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if len(repo.topics) != 0 {
		t.Errorf("got len(repo.topics) = %d; want 0", len(repo.topics))
	}
}

func TestAddSameTitleTwice(t *testing.T) {
	repo := NewInmemTopicRepository()

	repo.AddTopic("MyTopic")

	// We should not be able to add a task with same name twice
	err := repo.AddTopic("MyTopic")
	if err == nil {
		t.Errorf("got nil; want error")
	}
}

func TestRemove(t *testing.T) {
	repo := NewInmemTopicRepository()

	err := repo.AddTopic("MyTopic")
	if err != nil {
		t.Fatal(err)
	}
	repo.RemoveTopic("MyTopic")
	if len(repo.topics) != 0 {
		t.Errorf("got len(repo.topics) = %d; want 0", len(repo.topics))
	}
	if len(repo.topics) != 0 {
		t.Errorf("got len(repo.topicTitles) = %d; want 0", len(repo.topics))
	}
}

func TestGet(t *testing.T) {
	var titles = []string{
		"MyTopic1",
		"MyTopic2",
		"MyTopic3",
	}

	topicRepo := NewInmemTopicRepository()
	for _, title := range titles {
		topicRepo.AddTopic(title)
	}

	for _, title := range titles {
		topic, err := topicRepo.GetTopic(title)
		if err != nil {
			t.Error(err)
			continue
		}
		if topic.Title != title {
			t.Errorf("GetTopic(%s) want title %s; got %s",
				title, title, topic.Title)
		}
	}

	_, err := topicRepo.GetTopic("SomeTitle")
	if err == nil {
		t.Errorf("got nil; want err")
	}
}
