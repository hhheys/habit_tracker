package errors

import "errors"

var (
	ErrHabitNotFound      = errors.New("habit not found")
	ErrHabitAlreadyExists = errors.New("habit already linked")
	ErrHabitAlreadyAdded  = errors.New("habit already added")
)
