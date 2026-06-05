package tag

import (
	"context"
	"database/sql"
	"errors"
	"habit-tracker/internal/domain"

	"go.uber.org/zap"
)

type Repository struct {
	db *sql.DB

	log *zap.Logger
}

func NewRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, log: logger}
}

func (r *Repository) Create(ctx context.Context, tag *domain.HabitTag) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO tag (title) VALUES ($1) RETURNING id, title`, tag.Name).
		Scan(&tag.ID, &tag.Name)
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*domain.HabitTag, error) {
	var tag domain.HabitTag
	err := r.db.QueryRowContext(ctx, `SELECT id, title FROM tag WHERE id = $1`, id).Scan(&tag.ID, &tag.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrHabitNotFound
	}
	if err != nil {
		r.log.Error("failed to get tag by ID", zap.Error(err), zap.Uint("tag_id", id))
		return nil, err
	}
	return &tag, nil
}

func (r *Repository) Update(ctx context.Context, tag *domain.HabitTag) error {
	err := r.db.QueryRowContext(ctx, `
		UPDATE tag SET title = $1 WHERE id = $2 RETURNING id, title`,
		tag.Name, tag.ID,
	).Scan(&tag.ID, &tag.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrHabitNotFound
	}
	if err != nil {
		r.log.Error("failed to update tag", zap.Error(err), zap.Uint("tag_id", tag.ID))
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tag WHERE id = $1`, id)
	if err != nil {
		r.log.Error("failed to delete tag", zap.Error(err), zap.Uint("tag_id", id))
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		r.log.Error("failed to delete tag", zap.Error(err), zap.Uint("tag_id", id))
		return err
	}
	if n == 0 {
		return domain.ErrHabitNotFound
	}
	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*domain.HabitTag, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title FROM tag ORDER BY title, id`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			r.log.Error("failed to close tag rows", zap.Error(closeErr))
		}
	}(rows)
	tags := make([]*domain.HabitTag, 0)
	for rows.Next() {
		var tag domain.HabitTag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, rows.Err()
}
