package domain

import "time"

type DailyConfirmation struct {
	ID          uint
	UserHabitID uint
	ConfirmedAt time.Time
}

type Streak struct {
	UserHabitID      uint
	LongestStreak    int
	CurrentStreak    *int
	IsConfirmedToday bool
}

type HeatmapDay struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}
