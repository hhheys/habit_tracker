package habit

import "habit-tracker/internal/domain"

type ListHabitsParams struct {
	UserID    uint
	TagIDs    []uint
	Search    string
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}

type ListHabitsOutput struct {
	Habits []*domain.Habit
	Total  int64
	Limit  int
	Offset int
}

type CreateHabitInput struct {
	Title         string
	Description   string
	Tags          []uint
	ImageFilename string
}

type UpdateHabitInput struct {
	ID            uint
	Title         string
	Description   string
	ImageFilename string
	AddTagIDs     []uint
	RemoveTagIDs  []uint
}
