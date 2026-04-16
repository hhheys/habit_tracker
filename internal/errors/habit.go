package errors

import "errors"

var (
	ErrHabitNotFound         = errors.New("habit not found")
	ErrHabitAlreadyExists    = errors.New("habit already linked")
	ErrHabitAlreadyAdded     = errors.New("habit already added")
	ErrHabitNotAdded         = errors.New("habit not added") // Привычка не добавлена пользователем.
	ErrHabitAlreadyConfirmed = errors.New("habit already confirmed")
)
