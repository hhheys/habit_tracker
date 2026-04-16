package handler

import (
	"habit-tracker/internal/service"

	"go.uber.org/zap"
)

// Handler is the interface for the API handlers.
type Handler interface {
	UserHandler
	HabitHandler
	StreakHandler
}

// HandlerImpl is the implementation of the Handler interface.
type handlerImpl struct {
	log *zap.Logger

	UserHandler
	HabitHandler
	StreakHandler
}

// NewHandler returns a new Handler.
func NewHandler(service service.Service, log *zap.Logger) Handler {
	return &handlerImpl{
		log:           log,
		UserHandler:   NewUserHandler(log, service),
		HabitHandler:  NewHabitHandler(log, service),
		StreakHandler: NewStreakHandler(log, service),
	}
}
