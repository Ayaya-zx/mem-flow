package themes

import (
	"fmt"
	"sync"
	"time"
)

type InmemThemeStore struct {
	m           sync.Mutex
	themes      map[int]*Theme
	themeTitles map[string]struct{}
	nextId      int
}

func NewInmemThemeStore() *InmemThemeStore {
	return &InmemThemeStore{
		themes:      make(map[int]*Theme),
		themeTitles: make(map[string]struct{}),
	}
}

type ThemeTitleError string

func (e ThemeTitleError) Error() string {
	return string(e)
}

func (ts *InmemThemeStore) AddTheme(title string) error {
	if title == "" {
		return ThemeTitleError("theme's title is empty")
	}
	if _, ok := ts.themeTitles[title]; ok {
		return ThemeTitleError(fmt.Sprintf(
			"theme with title %s already exists",
			title,
		))
	}

	theme := new(Theme)
	theme.Title = title
	theme.Created = time.Now()
	theme.LastRepeated = theme.Created
	theme.NextRepeat = time.Now().Add(time.Duration(20 * time.Minute))

	ts.m.Lock()
	defer ts.m.Unlock()
	ts.themes[ts.nextId] = theme
	ts.nextId++
	ts.themeTitles[title] = struct{}{}
	return nil
}

func (ts *InmemThemeStore) RemoveTheme(id int) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	if theme, ok := ts.themes[id]; ok {
		delete(ts.themes, id)
		delete(ts.themeTitles, theme.Title)
	}
	return nil
}

func (ts *InmemThemeStore) ThemeRepeated(id int) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	t, ok := ts.themes[id]
	if !ok {
		return fmt.Errorf("theme with id %d does not exist", id)
	}
	t.Repeat()
	return nil
}

func (ts *InmemThemeStore) GetAllThemes() ([]*Theme, error) {
	ts.m.Lock()
	defer ts.m.Unlock()
	res := make([]*Theme, 0, len(ts.themes))
	for _, t := range ts.themes {
		res = append(res, t)
	}
	return res, nil
}

func (ts *InmemThemeStore) GetTheme(id int) (*Theme, error) {
	ts.m.Lock()
	defer ts.m.Unlock()
	t, ok := ts.themes[id]
	if !ok {
		return nil, fmt.Errorf("theme with id %d does not exist", id)
	}
	return t, nil
}
