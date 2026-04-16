package models

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
