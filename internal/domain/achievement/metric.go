package achievement

import (
	"time"

	"github.com/google/uuid"
)

// MetricScope represents a categorization or grouping of metrics within the system as a string type.
// It can be used to define the scope of a metric, such as "user" or "user_habit".
type MetricScope string

const (
	User      MetricScope = "user"
	UserHabit MetricScope = "user_habit"
)

type MetricKey string

const (
	TotalHabits   MetricKey = "total_habits"
	CurrentStreak MetricKey = "current_streak"
	LongestStreak MetricKey = "longest_streak"
)

type Metric struct {
	ID  uuid.UUID
	Key string
}

type UserMetric struct {
	ID        uuid.UUID
	UserID    uint
	MetricID  uuid.UUID
	Value     int
	UpdatedAt time.Time
}

type UserHabitMetric struct {
	ID          uuid.UUID
	UserHabitID uint
	MetricKey   MetricKey
	Value       int
	UpdatedAt   time.Time
}
