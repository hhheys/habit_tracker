package metric

import (
	"context"
	"encoding/json"
	domain "habit-tracker/internal/domain"
	"habit-tracker/internal/domain/achievement"
	"habit-tracker/internal/domain/events"
	"habit-tracker/internal/usecase/eventpublisher"
	"strconv"

	"go.uber.org/zap"
)

type Metric struct {
	metricRepository Repository
	txManager        TXManager
	outboxPublisher  EventPublisher
	userHabit        UserHabitRepository

	log *zap.Logger

	allowedEvents map[events.EventType]func(ctx context.Context, event events.Event) error
}

func NewEventService(metricRepository Repository, txManager TXManager, outboxPublisher EventPublisher, userHabit UserHabitRepository, logger *zap.Logger) *Metric {
	m := &Metric{
		metricRepository: metricRepository,
		txManager:        txManager,
		outboxPublisher:  outboxPublisher,
		userHabit:        userHabit,
		log:              logger,
	}
	m.initAllowedEvents()
	return m
}

func (m *Metric) initAllowedEvents() {
	m.allowedEvents = map[events.EventType]func(ctx context.Context, event events.Event) error{
		events.EventTypeHabitStreakConfirmed: m.processHabitStreakEvent,
		events.EventTypeUserHabitAdded:       m.processUserHabitEvent,
	}
}

func (m *Metric) ProcessEvent(ctx context.Context, event events.Event) error {
	if fn, ok := m.allowedEvents[event.EventType]; ok {
		return fn(ctx, event)
	}
	return nil
}

func (m *Metric) processHabitStreakEvent(ctx context.Context, event events.Event) error {
	var streak domain.Streak

	err := json.Unmarshal(event.Payload, &streak)
	if err != nil {
		m.log.Error("failed to unmarshal event payload", zap.Error(err))
		return err
	}

	if streak.CurrentStreak == nil {
		return nil
	}

	currentStreak := *streak.CurrentStreak

	err = m.txManager.WithTx(ctx, func(ctx context.Context) error {
		metric := achievement.UserHabitMetric{
			UserHabitID: streak.UserHabitID,
			MetricKey:   achievement.CurrentStreak,
			Value:       currentStreak,
		}
		updateErr := m.metricRepository.UpdateUserHabitMetric(ctx, &metric)
		if updateErr != nil {
			return updateErr
		}

		return eventpublisher.Publish(
			ctx,
			m.outboxPublisher,
			events.EventTypeUserHabitMetricUpdate,
			strconv.Itoa(int(metric.UserHabitID)),
			metric,
		)
	})

	return err
}

func (m *Metric) processUserHabitEvent(ctx context.Context, event events.Event) error {
	var userHabit domain.UserHabit

	err := json.Unmarshal(event.Payload, &userHabit)
	if err != nil {
		m.log.Error("failed to unmarshal event payload", zap.Error(err))
		return err
	}

	if userHabit.ID == 0 {
		return nil
	}

	err = m.txManager.WithTx(ctx, func(ctx context.Context) error {
		totalHabits, totalErr := m.userHabit.GetTotalUserHabits(ctx, userHabit.UserID)
		if totalErr != nil {
			return totalErr
		}

		metric := achievement.UserHabitMetric{
			UserHabitID: userHabit.ID,
			MetricKey:   achievement.TotalHabits,
			Value:       totalHabits,
		}
		updateErr := m.metricRepository.UpdateUserHabitMetric(ctx, &metric)
		if updateErr != nil {
			return updateErr
		}

		return eventpublisher.Publish(
			ctx,
			m.outboxPublisher,
			events.EventTypeUserHabitMetricUpdate,
			strconv.Itoa(int(metric.UserHabitID)),
			metric,
		)
	})

	return err
}
