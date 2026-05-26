package response

import (
	"fmt"
	"habit-tracker/internal/domain"
	userhabituc "habit-tracker/internal/usecase/userhabit"
	"time"
)

type AllHabitsResponse struct {
	Habits     []*HabitResponse   `json:"habits"`
	Pagination PaginationResponse `json:"pagination"`
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
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func NewAllHabitsResponse(habits []*domain.Habit, limit, offset int, total int64) *AllHabitsResponse {
	items := make([]*HabitResponse, len(habits))
	for i, habit := range habits {
		items[i] = NewHabitResponse(habit)
	}
	return &AllHabitsResponse{Habits: items, Pagination: NewPaginationResponse(total, limit, offset)}
}

func NewPaginationResponse(total int64, limit, offset int) PaginationResponse {
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}
	return PaginationResponse{Page: page, PageSize: limit, Total: total}
}

func NewHabitResponse(habit *domain.Habit) *HabitResponse {
	tags := make([]HabitTagResponse, len(habit.Tags))
	for i, tag := range habit.Tags {
		tags[i] = HabitTagResponse{ID: tag.ID, Name: tag.Name}
	}
	return &HabitResponse{
		ID:          habit.ID,
		Title:       habit.Title,
		Description: habit.Description,
		IsAdded:     habit.IsAdded,
		CreatedAt:   habit.CreatedAt,
		Tags:        tags,
		ImageURL:    fmt.Sprintf("/images/habits/%s", habit.ImageFilename),
	}
}

type HabitTagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewHabitTagResponse(tag *domain.HabitTag) *HabitTagResponse {
	return &HabitTagResponse{ID: tag.ID, Name: tag.Name}
}

type AllHabitTagsResponse struct {
	Tags []*HabitTagResponse `json:"tags"`
}

func NewHabitTagsResponse(tags []*domain.HabitTag) *AllHabitTagsResponse {
	items := make([]*HabitTagResponse, len(tags))
	for i, tag := range tags {
		items[i] = NewHabitTagResponse(tag)
	}
	return &AllHabitTagsResponse{Tags: items}
}

type UserHabitsResponse struct {
	Habits     []*UserHabitResponse `json:"habits"`
	Pagination PaginationResponse   `json:"pagination"`
}

type UserHabitResponse struct {
	ID      uint            `json:"id"`
	HabitID uint            `json:"habit_id"`
	AddedAt time.Time       `json:"added_at"`
	Streak  *StreakResponse `json:"streak,omitempty"`
}

func NewUserHabitsResponse(output *userhabituc.ListUserHabitsOutput) *UserHabitsResponse {
	items := make([]*UserHabitResponse, len(output.UserHabits))
	for i, item := range output.UserHabits {
		items[i] = &UserHabitResponse{
			ID:      item.Habit.ID,
			HabitID: item.Habit.HabitID,
			AddedAt: item.Habit.AddedAt,
			Streak:  NewStreakResponse(item.Streak),
		}
	}
	return &UserHabitsResponse{
		Habits:     items,
		Pagination: NewPaginationResponse(output.Total, output.Limit, output.Offset),
	}
}
