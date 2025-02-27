package entity

import (
	"time"
)

type Level int

type Topic struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Created      time.Time `json:"created"`
	LastRepeated time.Time `json:"lastRepeated"`
	NextRepeat   time.Time `json:"nextRepeat"`
	level        int
}

func NewTopic(id int, title string) *Topic {
	return &Topic{
		Id:           id,
		Title:        title,
		Created:      time.Now(),
		LastRepeated: time.Now(),
		NextRepeat:   time.Now().Add(20 * time.Minute),
		level:        0,
	}
}

func (t *Topic) Repeat() {
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
