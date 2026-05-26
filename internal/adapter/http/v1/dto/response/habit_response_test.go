package response

import (
	"testing"

	"habit-tracker/internal/domain"
	userhabituc "habit-tracker/internal/usecase/userhabit"
)

func TestNewUserHabitsResponseIncludesNestedHabit(t *testing.T) {
	output := &userhabituc.ListUserHabitsOutput{
		UserHabits: []*userhabituc.Output{{
			Habit: &domain.UserHabit{
				ID: 9,
				Habit: &domain.Habit{
					ID:            4,
					Title:         "Read",
					Description:   "Twenty pages",
					ImageFilename: "read.png",
				},
			},
		}},
	}

	response := NewUserHabitsResponse(output)
	got := response.Habits[0].Habit
	if got == nil {
		t.Fatal("nested habit is nil")
	}
	if got.ID != 4 || got.Title != "Read" || got.ImageURL != "/images/habits/read.png" {
		t.Fatalf("nested habit = %+v", got)
	}
}
