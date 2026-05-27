package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresTestDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "postgres:16",
		Env: map[string]string{
			"POSTGRES_DB":       "test_db",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "pass",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").SkipExternalCheck(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if host == "localhost" {
		host = "127.0.0.1"
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	connString := "postgres://user:pass@%s:%s/test_db"

	conn, err := sql.Open("postgres", fmt.Sprintf(connString, host, mappedPort.Port()))
	if err != nil {
		t.Fatal(err)
	}

	return conn
}
