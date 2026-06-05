package achievement

import (
	"context"
	"database/sql"
	"habit-tracker/internal/adapter/postgres/txmanager"
	domainachievement "habit-tracker/internal/domain/achievement"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) GetExpectedAchievements(ctx context.Context, userID uint, metricKey domainachievement.MetricKey) ([]*domainachievement.Achievement, error) {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	rows, err := executor.QueryContext(ctx, `
		SELECT
			a.id,
			a.code,
			a.title,
			COALESCE(a.description, ''),
			a.enabled,
			ac.metric_scope,
			m.id,
			m.metric_key,
			ac.operator,
			ac.value
		FROM achievement a
		JOIN achievement_condition ac ON ac.achievement_id = a.id
		JOIN metric m ON m.id = ac.metric_id
		WHERE a.enabled = TRUE
		  AND NOT EXISTS (
		      SELECT 1
		      FROM user_achievement ua
		      WHERE ua.user_id = $1
		        AND ua.achievement_id = a.id
		  )
		  AND EXISTS (
		      SELECT 1
		      FROM achievement_condition trigger_ac
		      JOIN metric trigger_m ON trigger_m.id = trigger_ac.metric_id
		      WHERE trigger_ac.achievement_id = a.id
		        AND trigger_m.metric_key = $2
		  )
		ORDER BY a.code, ac.id`,
		userID,
		metricKey,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.log.Error("failed to close achievement rows", zap.Error(closeErr))
		}
	}()

	achievementsByID := make(map[uuid.UUID]*domainachievement.Achievement)
	achievementOrder := make([]uuid.UUID, 0)

	for rows.Next() {
		var (
			a           domainachievement.Achievement
			condition   domainachievement.Condition
			metricID    uuid.UUID
			metricKeyDB string
		)

		if err := rows.Scan(
			&a.ID,
			&a.Code,
			&a.Title,
			&a.Description,
			&a.Enabled,
			&condition.MetricScope,
			&metricID,
			&metricKeyDB,
			&condition.Operator,
			&condition.TargetValue,
		); err != nil {
			return nil, err
		}

		existing, ok := achievementsByID[a.ID]
		if !ok {
			a.Conditions = make([]domainachievement.Condition, 0)
			achievementsByID[a.ID] = &a
			achievementOrder = append(achievementOrder, a.ID)
			existing = &a
		}

		condition.RequiredMetric = domainachievement.Metric{
			ID:  metricID,
			Key: metricKeyDB,
		}
		existing.Conditions = append(existing.Conditions, condition)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]*domainachievement.Achievement, 0, len(achievementOrder))
	for _, id := range achievementOrder {
		result = append(result, achievementsByID[id])
	}

	return result, nil
}

func (r *Repository) UnlockUserAchievementByCode(ctx context.Context, userID uint, code string) (*domainachievement.UserAchievement, error) {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	userAchievement := &domainachievement.UserAchievement{}
	err := executor.QueryRowContext(ctx, `
		WITH selected_achievement AS (
			SELECT id, code, title, COALESCE(description, '') AS description, enabled
			FROM achievement
			WHERE code = $2
		),
		unlocked AS (
			INSERT INTO user_achievement (user_id, achievement_id)
			SELECT $1, id
			FROM selected_achievement
			ON CONFLICT (user_id, achievement_id)
			DO UPDATE SET unlocked_at = user_achievement.unlocked_at
			RETURNING user_id, achievement_id, unlocked_at
		)
		SELECT
			u.user_id,
			u.unlocked_at,
			a.id,
			a.code,
			a.title,
			a.description,
			a.enabled
		FROM unlocked u
		JOIN selected_achievement a ON a.id = u.achievement_id`,
		userID,
		code,
	).Scan(
		&userAchievement.UserID,
		&userAchievement.UnlockedAt,
		&userAchievement.Achievement.ID,
		&userAchievement.Achievement.Code,
		&userAchievement.Achievement.Title,
		&userAchievement.Achievement.Description,
		&userAchievement.Achievement.Enabled,
	)
	if err != nil {
		return nil, err
	}

	return userAchievement, nil
}
