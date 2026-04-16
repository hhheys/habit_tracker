package service

import (
	"habit-tracker/internal/dto/response"
	"habit-tracker/internal/models"
	"habit-tracker/internal/repository/postgres"

	"go.uber.org/zap"
)

type StreakService interface {
	CreateDailyConfirmation(userID, habitID uint) error
}

type streakService struct {
	log        *zap.Logger
	repository postgres.Repository
}

func NewStreakService(logger *zap.Logger, repository postgres.Repository) StreakService {
	return &streakService{
		log:        logger,
		repository: repository,
	}
}

func (s *streakService) CreateDailyConfirmation(userID, habitID uint) error {
	userHabit, err := s.repository.GetUserHabit(userID, habitID)
	dailyConfirmation := models.DailyConfirmation{
		UserHabitID: userHabit.ID,
	}
	err = s.repository.CreateDailyConfirmation(&dailyConfirmation)
	if err != nil {
		return err
	}
	return nil
}

func (s *streakService) GetStreak(userID, habitID uint) (*response.StreakResponse, error) {
	userHabit, err := s.repository.GetUserHabit(userID, habitID)
	if err != nil {
		return nil, err
	}
	streak, err := s.repository.GetStreak(userHabit.ID)
	if err != nil {
		return nil, err
	}
	return response.NewStreakResponse(streak, userHabit.Habit.ID), nil
}
