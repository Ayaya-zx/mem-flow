package themes

import (
	"testing"
)

func TestAddAndGet(t *testing.T) {
	store := NewInmemThemeStore()
	nextId := store.nextId

	// Add task
	id, err := store.AddTheme("MyTheme")
	if err != nil {
		t.Fatal(err)
	}
	if len(store.themes) != 1 {
		t.Fatalf("got len(store.themes) = %d; want 1", len(store.themes))
	}
	if _, ok := store.themeTitles["MyTheme"]; !ok {
		t.Fatalf("got store.themeTitles[\"MyTheme\"] = false; want true")
	}
	if store.nextId == nextId {
		t.Fatalf("got store.nextId = %d; want %d", store.nextId, store.nextId+1)
	}

	// Get added task
	theme, err := store.GetTheme(id)
	if err != nil {
		t.Fatal(err)
	}
	if theme.Id != id {
		t.Errorf("got theme.Id = %d; id = %d", theme.Id, id)
	}
	if theme.Title != "MyTheme" {
		t.Errorf("got theme.Title = %s; want \"MyTheme\"", theme.Title)
	}
}

func TestAddEmpty(t *testing.T) {
	store := NewInmemThemeStore()
	id := store.nextId

	// Adding a task with an empty title is prohibited
	_, err := store.AddTheme("")
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if len(store.themes) != 0 {
		t.Errorf("got len(store.themes) = %d; want 0", len(store.themes))
	}
	if store.nextId != id {
		t.Errorf("got store.nextId = %d, want %d", store.nextId, id)
	}
}

func TestAddSameTitleTwice(t *testing.T) {
	store := NewInmemThemeStore()

	store.AddTheme("MyTheme")
	nextId := store.nextId

	// We should not be able to add a task with same name twice
	_, err := store.AddTheme("MyTheme")
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if store.nextId != nextId {
		t.Errorf("got store.nextId = %d, want %d", store.nextId, nextId)
	}
}

func TestRemove(t *testing.T) {
	store := NewInmemThemeStore()

	id, err := store.AddTheme("MyTheme")
	if err != nil {
		t.Fatal(err)
	}
	store.RemoveTheme(id)
	if len(store.themes) != 0 {
		t.Errorf("got len(store.themes) = %d; want 0", len(store.themes))
	}
	if len(store.themeTitles) != 0 {
		t.Errorf("got len(store.themeTitles) = %d; want 0", len(store.themeTitles))
	}
}

func TestRepeated(t *testing.T) {
	store := NewInmemThemeStore()

	id, err := store.AddTheme("MyTheme")
	if err != nil {
		t.Fatal(err)
	}

	theme, err := store.GetTheme(id)
	if err != nil {
		t.Fatal(err)
	}

	nextRep := theme.NextRepeat
	err = store.ThemeRepeated(id)
	if err != nil {
		t.Fatal(err)
	}
	if nextRep == theme.NextRepeat {
		t.Errorf("want theme.NextRepeat change; got %v", theme.NextRepeat)
	}

	err = store.ThemeRepeated(100)
	if err == nil {
		t.Errorf("got nil; want err")
	}
}

func TestGet(t *testing.T) {
	var tests = []struct {
		title  string
		wantId int
	}{
		{"MyTheme1", 1},
		{"MyTheme2", 2},
		{"MyTheme3", 3},
	}

	store := NewInmemThemeStore()
	for _, test := range tests {
		store.AddTheme(test.title)
	}

	for _, test := range tests {
		theme, err := store.GetTheme(test.wantId)
		if err != nil {
			t.Error(err)
			continue
		}
		if theme.Title != test.title {
			t.Errorf("on id = %d want title %s; got %s",
				test.wantId, test.title, theme.Title)
		}
	}

	_, err := store.GetTheme(100)
	if err == nil {
		t.Errorf("got nil; want err")
	}
}
