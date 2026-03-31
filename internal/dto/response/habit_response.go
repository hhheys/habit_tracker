package response

import (
	"habit-tracker/internal/models"
	"time"
)

type AllHabitsResponse struct {
	Habits     []*HabitResponse    `json:"habits"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type HabitResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsAdded     bool      `json:"is_added"`
	CreatedAt   time.Time `json:"created_at"`
	ImageURL    string    `json:"image_url"`
}

type PaginationResponse struct {
	Page     int   `json:"page,omitempty"`
	PageSize int   `json:"pageSize,omitempty"`
	Total    int64 `json:"total,omitempty"`
}

func NewAllHabitsResponse(habits []*models.Habit, page, pageSize int, total int64) *AllHabitsResponse {
	var habitResponses []*HabitResponse
	for _, habit := range habits {
		habitResponses = append(habitResponses, NewHabitResponse(habit))
	}
	return &AllHabitsResponse{
		Habits:     habitResponses,
		Pagination: NewPaginationResponse(total, page, pageSize),
	}
}

func NewPaginationResponse(total int64, page, pageSize int) *PaginationResponse {
	return &PaginationResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

func NewHabitResponse(habit *models.Habit) *HabitResponse {
	return &HabitResponse{
		ID:          habit.ID,
		Title:       habit.Title,
		Description: habit.Description,
		IsAdded:     habit.IsAdded,
		CreatedAt:   habit.CreatedAt,
		ImageURL:    habit.ImageFilename,
	}
}

type StreakResponse struct {
	HabitID          uint `json:"habit_id"`
	CurrentStreak    int  `json:"current_streak"`
	IsConfirmedToday bool `json:"is_confirmed_today"`
}
