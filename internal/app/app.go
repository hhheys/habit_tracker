package app

import (
	"context"
	"database/sql"
	"errors"
	"habit-tracker/config"
	"habit-tracker/config/logger"
	authadapter "habit-tracker/internal/adapter/auth"
	v1handler "habit-tracker/internal/adapter/http/v1/handler"
	v1router "habit-tracker/internal/adapter/http/v1/router"
	"habit-tracker/internal/adapter/postgres"
	authuc "habit-tracker/internal/usecase/auth"
	habituc "habit-tracker/internal/usecase/habit"
	streakuc "habit-tracker/internal/usecase/streak"
	taguc "habit-tracker/internal/usecase/tag"
	userhabituc "habit-tracker/internal/usecase/userhabit"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type App struct {
	DB      *sql.DB
	Config  config.Config
	Logger  *zap.Logger
	JWT     *authadapter.JwtService
	Handler v1handler.Handler
}

func NewApp(cfg config.Config) *App {
	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	db := postgres.CreateConnection(cfg)
	if db == nil {
		panic("could not connect to postgres")
	}
	repositories := postgres.NewRepositories(db, zapLogger)
	hasher := authadapter.NewHasher()
	jwt := authadapter.NewJWTService(cfg.JWTSecret, time.Hour)

	authService := authuc.NewService(
		repositories.Users,
		repositories.RefreshSessions,
		hasher,
		hasher,
		jwt,
		authadapter.NewRefreshTokenGenerator(),
	)
	habitService := habituc.NewService(repositories.Habits, repositories.Streaks)
	userHabitService := userhabituc.NewService(repositories.UserHabits, repositories.Streaks)
	streakService := streakuc.NewService(repositories.Streaks)
	tagService := taguc.NewService(repositories.Tags)

	return &App{
		DB:      db,
		Config:  cfg,
		Logger:  zapLogger,
		JWT:     jwt,
		Handler: v1handler.NewHandler(authService, habitService, userHabitService, &streakService, &tagService, zapLogger),
	}
}

func (app *App) SetupRouter() *gin.Engine {
	return v1router.NewRouter(app.Handler, app.JWT, app.Logger)
}

func (app *App) Run() {
	postgres.Migrate(app.DB)
	r := app.SetupRouter()
	logger.WithGin(app.Logger, r)
	srv := &http.Server{
		Addr:    net.JoinHostPort(app.Config.ServerHost, app.Config.ServerPort),
		Handler: r,
	}

	go func() {
		app.Logger.Info("starting server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
	_ = app.Logger.Sync()
}
