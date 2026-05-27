package response

import "habit-tracker/internal/domain"

type StreakResponse struct {
	UserHabitID      uint `json:"user_habit_id"`
	LongestStreak    int  `json:"longest_streak"`
	CurrentStreak    *int `json:"current_streak"`
	IsConfirmedToday bool `json:"is_confirmed_today"`
}

func NewStreakResponse(streak *domain.Streak) *StreakResponse {
	if streak == nil {
		return nil
	}
	return &StreakResponse{
		UserHabitID:      streak.UserHabitID,
		LongestStreak:    streak.LongestStreak,
		CurrentStreak:    streak.CurrentStreak,
		IsConfirmedToday: streak.IsConfirmedToday,
	}
}
