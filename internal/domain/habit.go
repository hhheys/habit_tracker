package domain

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

type UserHabit struct {
	ID      uint
	UserID  uint
	HabitID uint
	AddedAt time.Time
}

type HabitTag struct {
	ID   uint
	Name string
}
