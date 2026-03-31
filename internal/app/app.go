package app

import (
	"context"
	"database/sql"
	"habit-tracker/config"
	"habit-tracker/config/logger"
	"habit-tracker/internal/auth"
	"habit-tracker/internal/handler"
	"habit-tracker/internal/repository/postgres"
	"habit-tracker/internal/router"
	"habit-tracker/internal/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

// App - main app structure
type App struct {
	DB *sql.DB

	Config config.Config
	Logger *zap.Logger

	JWT auth.JWTService

	Handler handler.Handler
}

// NewApp - create new app instance
func NewApp(config config.Config) *App {
	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	dbConn := postgres.CreateConnection(config)
	appRepository := postgres.NewRepository(dbConn, zapLogger)
	jwtService := auth.NewJWTService(config.JWTSecret, time.Minute*60)
	appService := service.NewService(zapLogger, appRepository, jwtService)
	appHandler := handler.NewHandler(appService, zapLogger)

	return &App{
		DB: dbConn,

		Config: config,
		Logger: zapLogger,

		JWT: jwtService,

		Handler: appHandler,
	}
}

// SetupRouter setup router for app
func (app *App) SetupRouter() *gin.Engine {
	return router.NewRouter(app.Handler, app.JWT)
}

// Run runs the application
func (app *App) Run() {
	postgres.Migrate(app.DB)

	r := app.SetupRouter()

	srv := &http.Server{
		Addr:    net.JoinHostPort(app.Config.ServerHost, app.Config.ServerPort),
		Handler: r,
	}

	logger.WithGin(app.Logger, r)

	app.Logger.Info("Starting server", zap.String("host", app.Config.ServerHost), zap.String("port", app.Config.ServerPort))

	if err := r.Run(net.JoinHostPort(app.Config.ServerHost, app.Config.ServerPort)); err != nil {
		app.Logger.Error("Server failed", zap.Error(err))
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		_ = app.Logger.Sync()
		cancel()
		log.Printf("Server forced to shutdown: %v", err)
		return
	}

	log.Println("Server exiting gracefully")
}
