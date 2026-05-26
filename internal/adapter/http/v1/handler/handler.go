package handler

import (
	authuc "habit-tracker/internal/usecase/auth"
	habituc "habit-tracker/internal/usecase/habit"
	streakuc "habit-tracker/internal/usecase/streak"
	taguc "habit-tracker/internal/usecase/tag"
	userhabituc "habit-tracker/internal/usecase/userhabit"

	"go.uber.org/zap"
)

type Handler interface {
	AuthHandler
	HabitHandler
	StreakHandler
}

type handlerImpl struct {
	AuthHandler
	HabitHandler
	StreakHandler
}

func NewHandler(
	auth *authuc.Service,
	habits *habituc.Service,
	userHabits *userhabituc.Service,
	streaks *streakuc.Service,
	tags *taguc.Service,
	log *zap.Logger,
) Handler {
	return &handlerImpl{
		AuthHandler:   NewAuthHandler(auth, log),
		HabitHandler:  NewHabitHandler(habits, userHabits, tags),
		StreakHandler: NewStreakHandler(streaks),
	}
}
