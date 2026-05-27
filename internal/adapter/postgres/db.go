package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"habit-tracker/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// CreateConnection creates a new connection to the database
func CreateConnection(config config.Config) *sql.DB {
	attempts := 0

	for attempts < 3 {
		attempts++
		conn, err := sql.Open("postgres", config.PGConnString)
		if err != nil {
			log.Println("Failed to open connection, retrying...")
			continue
		}

		if err := conn.Ping(); err != nil {
			log.Println("Cannot connect to DB:", err)
			continue
		}
		return conn
	}
	return nil
}

// Migrate applies migrations to the database
func Migrate(conn *sql.DB) {
	MigrateFrom(conn, "file://./migrations")
}

// MigrateFrom applies migrations from the given source URL to the database.
func MigrateFrom(conn *sql.DB, sourceURL string) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		panic(fmt.Sprintf("Couldn't create migration driver: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		panic(fmt.Sprintf("Coulnd't create migrator: %v", err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(fmt.Sprintf("Couldn't apply migrations: %v", err))
	}
}
