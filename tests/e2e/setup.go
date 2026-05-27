package e2e

import (
	"database/sql"
	"habit-tracker/config"
	"habit-tracker/internal/adapter/postgres"
	"habit-tracker/internal/app"
	"habit-tracker/tests/e2e/testdb"
	"testing"

	"github.com/gin-gonic/gin"
)

type TestApp struct {
	Router *gin.Engine
	DB     *sql.DB
}

func setupTestApp(t *testing.T) *TestApp {
	db := testdb.NewPostgresTestDB(t)

	t.Cleanup(func() {
		_ = db.Close()
	})

	gin.SetMode(gin.TestMode)

	appConfig := config.Config{
		JWTSecret: "test-jwt-secret",
	}

	postgres.MigrateFrom(db, "file://../../migrations")

	appInstance := app.NewAppWithDB(appConfig, db)

	router := appInstance.SetupRouter()

	return &TestApp{
		Router: router,
		DB:     db,
	}
}
