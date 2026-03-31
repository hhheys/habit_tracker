package config

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Config contains configuration for the application.
type Config struct {
	ServerHost string
	ServerPort string

	PGConnString string

	JWTSecret string
}

// NewConfig creates a new Config instance.
func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env files found")
	}
	return Config{
		ServerHost:   getEnvOrFatal("SERVER_HOST"),
		ServerPort:   getEnvOrFatal("SERVER_PORT"),
		PGConnString: GetPostgresURL(),
		JWTSecret:    getEnvOrFatal("JWT_SECRET"),
	}
}

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Empty environment variable: %s", key)
	}
	return value
}

// GetPostgresURL returns the postgres connection string in the format: postgresql://
func GetPostgresURL() string {
	host := getEnvOrFatal("POSTGRES_HOST")
	port := getEnvOrFatal("POSTGRES_PORT")
	user := getEnvOrFatal("POSTGRES_USER")
	password := getEnvOrFatal("POSTGRES_PASSWORD")
	dbname := getEnvOrFatal("POSTGRES_DB")
	sslMode := getEnvOrFatal("POSTGRES_SSL_MODE")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		dbname,
		sslMode,
	)
}
