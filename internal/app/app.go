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
	kafkaadapter "habit-tracker/internal/adapter/kafka"
	"habit-tracker/internal/adapter/postgres"
	"habit-tracker/internal/adapter/postgres/txmanager"
	authuc "habit-tracker/internal/usecase/auth"
	habituc "habit-tracker/internal/usecase/habit"
	metricuc "habit-tracker/internal/usecase/metric"
	streakuc "habit-tracker/internal/usecase/streak"
	taguc "habit-tracker/internal/usecase/tag"
	userhabituc "habit-tracker/internal/usecase/userhabit"
	outboxworker "habit-tracker/internal/worker/outbox"
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

	KafkaProducer   *kafkaadapter.Producer
	KafkaConsumer   *kafkaadapter.Consumer
	OutboxPublisher outboxworker.EventPublisher
	MetricService   *metricuc.Metric
}

func NewApp(cfg config.Config) *App {
	db := postgres.CreateConnection(cfg)
	if db == nil {
		panic("could not connect to postgres")
	}

	return NewAppWithDB(cfg, db)
}

func NewAppWithDB(cfg config.Config, db *sql.DB) *App {
	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	repositories := postgres.NewRepositories(db, zapLogger)
	hasher := authadapter.NewHasher()
	jwt := authadapter.NewJWTService(cfg.JWTSecret, time.Hour)
	txManager := txmanager.NewTXManager(db)
	kafkaProducer := kafkaadapter.NewProducer(cfg.KafkaBrokers, zapLogger)
	kafkaConsumer := kafkaadapter.NewConsumer(cfg.KafkaBrokers, zapLogger)
	outboxPublisher := outboxworker.NewEventPublisher(repositories.Outbox, kafkaProducer, cfg.KafkaOutboxTopic, zapLogger)

	authService := authuc.NewService(
		repositories.Users,
		repositories.RefreshSessions,
		hasher,
		hasher,
		jwt,
		authadapter.NewRefreshTokenGenerator(),
	)
	habitService := habituc.NewService(repositories.Habits, repositories.Streaks)
	userHabitService := userhabituc.NewService(repositories.UserHabits, repositories.Streaks, repositories.Outbox, txManager)
	streakService := streakuc.NewService(repositories.Streaks, repositories.Outbox, txManager)
	tagService := taguc.NewService(repositories.Tags)
	metricService := metricuc.NewEventService(repositories.Metrics, txManager, repositories.Outbox, repositories.UserHabits, zapLogger)

	return &App{
		DB:      db,
		Config:  cfg,
		Logger:  zapLogger,
		JWT:     jwt,
		Handler: v1handler.NewHandler(authService, habitService, userHabitService, &streakService, &tagService, zapLogger),

		KafkaProducer:   kafkaProducer,
		KafkaConsumer:   kafkaConsumer,
		OutboxPublisher: outboxPublisher,
		MetricService:   metricService,
	}
}

func (app *App) SetupRouter() *gin.Engine {
	return v1router.NewRouter(app.Handler, app.JWT, app.Logger)
}

func (app *App) Run() {
	postgres.Migrate(app.DB)
	workerCtx, stopWorker := context.WithCancel(context.Background())
	defer stopWorker()

	go func() {
		app.Logger.Info(
			"starting outbox publisher",
			zap.Strings("brokers", app.Config.KafkaBrokers),
			zap.String("topic", app.Config.KafkaOutboxTopic),
			zap.Duration("interval", app.Config.OutboxPublishInterval),
		)
		app.OutboxPublisher.Run(workerCtx, app.Config.OutboxPublishInterval)
	}()

	go func() {
		app.Logger.Info(
			"starting metric event consumer",
			zap.Strings("brokers", app.Config.KafkaBrokers),
			zap.String("topic", app.Config.KafkaOutboxTopic),
		)
		if err := app.KafkaConsumer.ConsumeEvents(
			workerCtx,
			app.Config.KafkaOutboxTopic,
			"metric-service",
			app.MetricService.ProcessEvent,
		); err != nil && !errors.Is(err, context.Canceled) {
			app.Logger.Error("metric event consumer stopped", zap.Error(err))
		}
	}()

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

	stopWorker()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
	if err := app.KafkaProducer.Close(); err != nil {
		app.Logger.Error("failed to close kafka producer", zap.Error(err))
	}
	if err := app.KafkaConsumer.Close(); err != nil {
		app.Logger.Error("failed to close kafka consumer", zap.Error(err))
	}
	_ = app.Logger.Sync()
}
