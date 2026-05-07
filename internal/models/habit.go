package models

import (
	"time"
)

type Habit struct {
	ID            uint
	Title         string
	Description   string
	IsAdded       bool
	CreatedAt     time.Time
	Tags          []*HabitTag
	ImageFilename string
}

func NewHabit(title, description string) *Habit {
	return &Habit{
		Title:       title,
		Description: description,
	}
}

type UserHabit struct {
	ID      uint
	UserID  uint
	Habit   Habit
	AddedAt time.Time
	Streak  *Streak
}

type HabitTag struct {
	ID    uint
	Title string
}
