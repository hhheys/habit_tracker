package achievement

import (
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	ID          uuid.UUID
	Code        string
	Title       string
	Description string
	Enabled     bool
	Conditions  []Condition
}

type Condition struct {
	MetricScope    MetricScope
	RequiredMetric Metric
	Operator       string
	TargetValue    int
}

type UserAchievement struct {
	UserID      uint
	Achievement Achievement
	UnlockedAt  time.Time
}
