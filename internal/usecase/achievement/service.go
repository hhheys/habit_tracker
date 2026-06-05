package achievement

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"habit-tracker/internal/domain/achievement"
	"habit-tracker/internal/domain/events"
	"habit-tracker/internal/usecase/eventpublisher"
	"strconv"

	"go.uber.org/zap"
)

type Service struct {
	userAchievement UserAchievementRepository
	userMetric      UserMetricRepository
	userHabit       UserHabitRepository
	txManager       TXManager
	eventPublisher  EventPublisher

	log *zap.Logger
}

func NewService(log *zap.Logger, userAchievement UserAchievementRepository, userMetric UserMetricRepository, userHabit UserHabitRepository, txManager TXManager, eventPublisher EventPublisher) *Service {
	return &Service{userAchievement: userAchievement, userMetric: userMetric, userHabit: userHabit, txManager: txManager, eventPublisher: eventPublisher, log: log}
}

func (s *Service) ProcessEvent(ctx context.Context, event events.Event) error {
	if event.EventType != events.EventTypeUserHabitMetricUpdate {
		return nil
	}

	var metric achievement.UserHabitMetric
	if err := json.Unmarshal(event.Payload, &metric); err != nil {
		s.log.Error("failed to unmarshal user habit metric event payload", zap.Error(err))
		return err
	}

	if metric.UserHabitID == 0 {
		return nil
	}

	userID, err := s.userHabit.GetUserIDByUserHabitID(ctx, metric.UserHabitID)
	if err != nil {
		s.log.Error("failed to get user id by user habit id", zap.Error(err), zap.Uint("user_habit_id", metric.UserHabitID))
		return err
	}

	return s.processAchievements(ctx, userID, metric.UserHabitID, metric.MetricKey)
}

func (s *Service) ProcessAchievements(
	ctx context.Context,
	userID uint,
	metricKey achievement.MetricKey,
) error {
	return s.processAchievements(ctx, userID, 0, metricKey)
}

func (s *Service) ListUserAchievements(ctx context.Context, params ListUserAchievementsParams) (*ListUserAchievementsOutput, error) {
	achievements, total, err := s.userAchievement.ListUserAchievements(ctx, params.UserID, params.Limit, params.Offset)
	if err != nil {
		s.log.Error("failed to list user achievements", zap.Error(err), zap.Uint("user_id", params.UserID))
		return nil, err
	}

	return &ListUserAchievementsOutput{
		Achievements: achievements,
		Limit:        params.Limit,
		Offset:       params.Offset,
		Total:        total,
	}, nil
}

func (s *Service) processAchievements(
	ctx context.Context,
	userID uint,
	userHabitID uint,
	metricKey achievement.MetricKey,
) error {
	// Collecting user metrics as needed
	userMetrics := make(map[achievement.MetricKey]int)
	userHabitMetrics := make(map[achievement.MetricKey]int)

	expectedAchievements, err := s.userAchievement.GetExpectedAchievements(ctx, userID, metricKey)
	if err != nil {
		s.log.Error("failed to get expected achievements", zap.Error(err))
		return err
	}

	for _, a := range expectedAchievements {
		passedConditions := 0
		for _, condition := range a.Conditions {
			conditionMetricKey := achievement.MetricKey(condition.RequiredMetric.Key)
			v, metricErr := s.getConditionMetricValue(ctx, userID, userHabitID, condition.MetricScope, conditionMetricKey, userMetrics, userHabitMetrics)
			if errors.Is(metricErr, sql.ErrNoRows) {
				break
			}
			if metricErr != nil {
				s.log.Error(
					"failed to get condition metric",
					zap.Error(metricErr),
					zap.String("metric_key", string(conditionMetricKey)),
					zap.String("metric_scope", string(condition.MetricScope)),
				)
				return metricErr
			}

			if s.processConditionOperator(condition.Operator, condition.TargetValue, v) {
				passedConditions++
			}
		}

		if len(a.Conditions) == 0 {
			continue
		}

		if passedConditions == len(a.Conditions) {
			txErr := s.txManager.WithTx(
				ctx,
				func(customContext context.Context) error {
					ua, unlockErr := s.userAchievement.UnlockUserAchievementByCode(customContext, userID, a.Code)
					if unlockErr != nil {
						s.log.Error("failed to unlock user achievement", zap.Error(unlockErr), zap.String("achievement_code", a.Code))
						return unlockErr
					}

					return eventpublisher.Publish(
						customContext,
						s.eventPublisher,
						events.EventTypeUserAchievementUnlocked,
						strconv.Itoa(int(userID)),
						ua,
					)
				},
			)

			if txErr != nil {
				s.log.Error("failed to unlock user achievement", zap.Error(txErr))
				return txErr
			}
		}
	}
	return nil
}

func (s *Service) getConditionMetricValue(
	ctx context.Context,
	userID uint,
	userHabitID uint,
	scope achievement.MetricScope,
	metricKey achievement.MetricKey,
	userMetrics map[achievement.MetricKey]int,
	userHabitMetrics map[achievement.MetricKey]int,
) (int, error) {
	switch scope {
	case achievement.User:
		v, ok := userMetrics[metricKey]
		if !ok {
			uMetric, err := s.userMetric.GetUserMetricByKeyAndUserID(ctx, userID, string(metricKey))
			if err != nil {
				return 0, err
			}
			userMetrics[metricKey] = uMetric.Value
			v = uMetric.Value
		}
		return v, nil
	case achievement.UserHabit:
		if userHabitID == 0 {
			return 0, sql.ErrNoRows
		}

		v, ok := userHabitMetrics[metricKey]
		if !ok {
			uMetric, err := s.userMetric.GetUserHabitMetricByKeyAndUserHabitID(ctx, userHabitID, string(metricKey))
			if err != nil {
				return 0, err
			}
			userHabitMetrics[metricKey] = uMetric.Value
			v = uMetric.Value
		}
		return v, nil
	default:
		v, ok := userMetrics[metricKey]
		if !ok {
			uMetric, err := s.userMetric.GetUserMetricByKeyAndUserID(ctx, userID, string(metricKey))
			if errors.Is(err, sql.ErrNoRows) && userHabitID != 0 {
				uHabitMetric, habitMetricErr := s.userMetric.GetUserHabitMetricByKeyAndUserHabitID(ctx, userHabitID, string(metricKey))
				if habitMetricErr != nil {
					return 0, habitMetricErr
				}
				userHabitMetrics[metricKey] = uHabitMetric.Value
				return uHabitMetric.Value, nil
			}
			if err != nil {
				return 0, err
			}
			userMetrics[metricKey] = uMetric.Value
			v = uMetric.Value
		}
		return v, nil
	}
}

func (s *Service) processConditionOperator(operator string, targetValue, actualValue int) bool {
	switch operator {
	case "==":
		return targetValue == actualValue
	case ">=":
		return actualValue >= targetValue
	case "<=":
		return actualValue <= targetValue
	default:
		return false
	}
}
