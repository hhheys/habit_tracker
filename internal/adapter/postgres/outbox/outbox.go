package outbox

import (
	"context"
	"database/sql"
	"habit-tracker/internal/adapter/postgres/txmanager"
	"habit-tracker/internal/domain/events"
	"time"

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

func (r *Repository) Publish(ctx context.Context, event *events.Event) error {
	executor := txmanager.ExecutorFromContext(ctx, r.db)
	err := executor.QueryRowContext(
		ctx,
		`INSERT INTO outbox_event (
                          occurred_at, 
                          event_type, 
                          event_type_version, 
                          partition_key, 
                          payload) 
				VALUES ($1, $2, $3, $4,$5) 
				RETURNING id`,
		event.OccurredAt,
		event.EventType,
		event.EventTypeVersion,
		event.PartitionKey,
		event.Payload,
	).Scan(&event.EventID)
	if err != nil {
		r.log.Error("failed to publish event to outbox", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) GetCreated(ctx context.Context, limit int) ([]*events.Event, error) {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	rows, err := executor.QueryContext(
		ctx,
		`SELECT 
			id,
			occurred_at,
			event_type,
			event_type_version,
			partition_key,
			payload,
			attempt_count
		FROM outbox_event
		WHERE status = 'created'
		  AND next_attempt_at <= NOW()
		ORDER BY created_at
		LIMIT $1`,
		limit,
	)
	if err != nil {
		r.log.Error("failed to get created outbox events", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	result := make([]*events.Event, 0)

	for rows.Next() {
		var event events.Event

		err := rows.Scan(
			&event.EventID,
			&event.OccurredAt,
			&event.EventType,
			&event.EventTypeVersion,
			&event.PartitionKey,
			&event.Payload,
			&event.AttemptCount,
		)
		if err != nil {
			r.log.Error("failed to scan outbox event", zap.Error(err))
			return nil, err
		}

		result = append(result, &event)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate outbox events", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) MarkDead(ctx context.Context, id uuid.UUID) error {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	_, err := executor.ExecContext(
		ctx,
		`UPDATE outbox_event
		 SET 
			status = 'dead',
			updated_at = NOW()
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		r.log.Error("failed to mark outbox event as dead", zap.Error(err), zap.String("event_id", id.String()))
		return err
	}

	return nil
}

func (r *Repository) IncrementAttemptCountAndNextTime(ctx context.Context, id uuid.UUID, nextTime time.Time) error {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	_, err := executor.ExecContext(
		ctx,
		`UPDATE outbox_event
		 SET 
			attempt_count = attempt_count + 1,
			next_attempt_at = $2,
			updated_at = NOW()
		 WHERE id = $1`,
		id,
		nextTime,
	)
	if err != nil {
		r.log.Error("failed to increment outbox event attempt count", zap.Error(err), zap.String("event_id", id.String()))
		return err
	}

	return nil
}

func (r *Repository) MarkSent(ctx context.Context, id uuid.UUID) error {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	_, err := executor.ExecContext(
		ctx,
		`UPDATE outbox_event
		 SET 
			status = 'sent',
			updated_at = NOW()
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		r.log.Error("failed to mark outbox event as sent", zap.Error(err), zap.String("event_id", id.String()))
		return err
	}

	return nil
}
