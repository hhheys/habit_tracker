package postgres

import (
	"database/sql"
	pgachievement "habit-tracker/internal/adapter/postgres/achievement"
	pghabit "habit-tracker/internal/adapter/postgres/habit"
	pgmetric "habit-tracker/internal/adapter/postgres/metric"
	pgoutbox "habit-tracker/internal/adapter/postgres/outbox"
	pgsession "habit-tracker/internal/adapter/postgres/session"
	pgstreak "habit-tracker/internal/adapter/postgres/streak"
	pgtag "habit-tracker/internal/adapter/postgres/tag"
	pguser "habit-tracker/internal/adapter/postgres/user"
	pguserhabit "habit-tracker/internal/adapter/postgres/userhabit"

	"go.uber.org/zap"
)

type Repositories struct {
	Users           *pguser.Repository
	RefreshSessions *pgsession.Repository
	Habits          *pghabit.Repository
	UserHabits      *pguserhabit.Repository
	Streaks         *pgstreak.Repository
	Tags            *pgtag.Repository
	Outbox          *pgoutbox.Repository
	Metrics         *pgmetric.Repository
	Achievements    *pgachievement.Repository
}

func NewRepositories(db *sql.DB, log *zap.Logger) *Repositories {
	return &Repositories{
		Users:           pguser.NewRepository(db, log),
		RefreshSessions: pgsession.NewRepository(db, log),
		Habits:          pghabit.NewRepository(db, log),
		UserHabits:      pguserhabit.NewRepository(db, log),
		Streaks:         pgstreak.NewRepository(db, log),
		Tags:            pgtag.NewRepository(db, log),
		Outbox:          pgoutbox.NewRepository(db, log),
		Metrics:         pgmetric.NewRepository(db, log),
		Achievements:    pgachievement.NewRepository(db, log),
	}
}
