package response

import "habit-tracker/internal/models"

type StreakResponse struct {
	HabitID          uint `json:"habit_id"`
	CurrentStreak    *int `json:"current_streak"`
	IsConfirmedToday bool `json:"is_confirmed_today"`
}

func NewStreakResponse(streak *models.Streak, habitID uint) *StreakResponse {
	return &StreakResponse{
		HabitID:          habitID,
		CurrentStreak:    streak.CurrentStreak,
		IsConfirmedToday: streak.IsConfirmedToday,
	}
}
