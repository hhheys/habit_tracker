package response

import (
	"fmt"
	"habit-tracker/internal/models"
	"time"
)

type AllHabitsResponse struct {
	Habits     []*HabitResponse    `json:"habits"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type UserHabitsResponse struct {
	Habits []*UserHabitResponse `json:"habits"`
	//Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type UserHabitResponse struct {
	Habit   *HabitResponse  `json:"habit"`
	AddedAt time.Time       `json:"added_at"`
	Streak  *StreakResponse `json:"streak,omitempty"`
}

func NewUserHabitResponse(habit *models.UserHabit) *UserHabitResponse {
	return &UserHabitResponse{
		Habit: &HabitResponse{
			ID:          habit.Habit.ID,
			Title:       habit.Habit.Title,
			Description: habit.Habit.Description,
			IsAdded:     habit.Habit.IsAdded,
			ImageURL:    fmt.Sprintf("%s%s", "/images/habits/", habit.Habit.ImageFilename),
		},
		AddedAt: habit.AddedAt,
		Streak:  NewStreakResponse(habit.Streak, habit.ID),
	}
}

func NewUserHabitsResponse(habits []*models.UserHabit) *UserHabitsResponse {
	res := make([]*UserHabitResponse, len(habits))
	for i, habit := range habits {
		res[i] = NewUserHabitResponse(habit)
	}
	return &UserHabitsResponse{res}
}

type HabitResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	IsAdded     bool               `json:"is_added"`
	CreatedAt   time.Time          `json:"created_at"`
	ImageURL    string             `json:"image_url"`
	Tags        []HabitTagResponse `json:"tags"`
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
	tags := make([]HabitTagResponse, len(habit.Tags))
	for i, habitTag := range habit.Tags {
		tags[i] = HabitTagResponse{
			ID:   habitTag.ID,
			Name: habitTag.Title,
		}
	}
	return &HabitResponse{
		ID:          habit.ID,
		Title:       habit.Title,
		Description: habit.Description,
		IsAdded:     habit.IsAdded,
		CreatedAt:   habit.CreatedAt,
		Tags:        tags,
		ImageURL:    fmt.Sprintf("%s%s", "/images/habits/", habit.ImageFilename),
	}
}

type HabitTagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AllHabitTagsResponse struct {
	Tags       []*HabitTagResponse `json:"tags"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

func NewHabitTagResponse(tag *models.HabitTag) *HabitTagResponse {
	return &HabitTagResponse{
		ID:   tag.ID,
		Name: tag.Title,
	}
}

func NewHabitTagsResponse(tags []*models.HabitTag, pagination *PaginationResponse) *AllHabitTagsResponse {
	res := make([]*HabitTagResponse, len(tags))
	for i, tag := range tags {
		res[i] = NewHabitTagResponse(tag)
	}
	return &AllHabitTagsResponse{
		Tags:       res,
		Pagination: pagination,
	}
}

type EditHabitResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
