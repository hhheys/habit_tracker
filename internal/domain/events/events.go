package events

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

type Event struct {
	EventID          uuid.UUID
	OccurredAt       time.Time
	EventType        EventType
	EventTypeVersion int
	AttemptCount     int
	PartitionKey     string
	Payload          []byte
}

func (e Event) GetEventID() uuid.UUID {
	return e.EventID
}

const (
	EventTypeUserHabitAdded       EventType = "user_habit.added.v1"
	EventTypeUserHabitRemoved     EventType = "user_habit.removed.v1"
	EventTypeHabitStreakConfirmed EventType = "habit_streak.confirmed.v1"
)
