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
	Habit   *Habit
}

type HabitTag struct {
	ID   uint
	Name string
}

type HabitListFilter struct {
	UserID    uint
	TagIDs    []uint
	Search    string
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}

type UserHabitListFilter struct {
	UserID    uint
	Search    string
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}
