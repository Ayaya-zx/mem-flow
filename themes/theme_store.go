package themes

import "time"

type Level int

type Theme struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Created      time.Time `json:"created"`
	LastRepeated time.Time `json:"lastRepeated"`
	NextRepeat   time.Time `json:"nextRepeat"`
	level        int
}

func (t *Theme) Repeat() {
	t.LastRepeated = time.Now()
	switch t.level {
	case 0: // 8 hours
		t.NextRepeat = time.Now().Add(time.Duration(8 * time.Hour))
	case 1: // 24 hours
		t.NextRepeat = time.Now().Add(time.Duration(24 * time.Hour))
	case 2: // 2 days
		t.NextRepeat = time.Now().Add(time.Duration(2 * 24 * time.Hour))
	case 3: // 1 week
		t.NextRepeat = time.Now().Add(time.Duration(7 * 24 * time.Hour))
	default: // 1 month
		t.NextRepeat = time.Now().Add(time.Duration(30 * 24 * time.Hour))
	}
	if t.level < 4 {
		t.level++
	}
}

type ThemeStore interface {
	AddTheme(string) error
	RemoveTheme(int) error
	ThemeRepeated(int) error
	GetAllThemes() ([]*Theme, error)
	GetTheme(int) (*Theme, error)
}
