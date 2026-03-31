package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"habit-tracker/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// CreateConnection creates a new connection to the database
func CreateConnection(config config.Config) *sql.DB {
	conn, err := sql.Open("postgres", config.PGConnString)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	return conn
}

// Migrate applies migrations to the database
func Migrate(conn *sql.DB) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		panic(fmt.Sprintf("Couldn't create migration driver: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
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
