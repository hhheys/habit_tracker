package userhabit

import "habit-tracker/internal/domain"

type Output struct {
	Habit  *domain.UserHabit
	Streak *domain.Streak
}

type ListUserHabitsOutput struct {
	UserHabits []*Output
	Total      int64
	Limit      int
	Offset     int
}

type ListUserHabitsParams struct {
	UserID    uint
	Search    string
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}

type AddUserHabitInput struct {
	UserID  uint
	HabitID uint
}
