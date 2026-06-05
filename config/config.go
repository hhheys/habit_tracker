package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config contains configuration for the application.
type Config struct {
	ServerHost string
	ServerPort string

	PGConnString string

	JWTSecret string

	KafkaBrokers          []string
	KafkaOutboxTopic      string
	OutboxPublishInterval time.Duration
}

// NewConfig creates a new Config instance.
func NewConfig() Config {
	_ = godotenv.Overload()
	return Config{
		ServerHost:            getEnvOrFatal("SERVER_HOST"),
		ServerPort:            getEnvOrFatal("SERVER_PORT"),
		PGConnString:          GetPostgresURL(),
		JWTSecret:             getEnvOrFatal("JWT_SECRET"),
		KafkaBrokers:          getCSVEnvOrDefault("KAFKA_BROKERS", "localhost:9092"),
		KafkaOutboxTopic:      getEnvOrDefault("KAFKA_OUTBOX_TOPIC", "habit.events.v1"),
		OutboxPublishInterval: getDurationEnvOrDefault("OUTBOX_PUBLISH_INTERVAL", 5*time.Second),
	}
}

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Empty environment variable: %s", key)
	}
	return value
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getCSVEnvOrDefault(key, fallback string) []string {
	rawValue := getEnvOrDefault(key, fallback)
	parts := strings.Split(rawValue, ",")
	values := make([]string, 0, len(parts))

	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}

	if len(values) == 0 {
		return []string{fallback}
	}

	return values
}

func getDurationEnvOrDefault(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid duration in environment variable %s: %v", key, err)
	}

	return duration
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
