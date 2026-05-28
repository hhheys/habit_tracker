APP_NAME := habit-tracker
SERVER_PKG := ./cmd/server
BIN_DIR := bin
BIN := $(BIN_DIR)/server
COMPOSE := docker compose
DB_CONTAINER := habit-tracker-db
DB_USER ?= postgres
DB_NAME ?= habit_tracker

.PHONY: help run build test fmt tidy docker-up docker-rebuild docker-down docker-logs logs db-shell migrate-status clean

help:
	@echo "Available targets:"
	@echo "  make run             Run app locally"
	@echo "  make build           Build server binary into $(BIN)"
	@echo "  make test            Run all Go tests"
	@echo "  make fmt             Format Go files"
	@echo "  make tidy            Run go mod tidy"
	@echo "  make docker-up       Start docker compose services"
	@echo "  make docker-rebuild  Rebuild and start app service"
	@echo "  make docker-down     Stop docker compose services"
	@echo "  make logs            Tail app logs"
	@echo "  make docker-logs     Tail all compose logs"
	@echo "  make db-shell        Open psql in Postgres container"
	@echo "  make migrate-status  Show migration version and outbox table"
	@echo "  make clean           Remove local build artifacts"

run:
	go run $(SERVER_PKG)

build:
	go build -o $(BIN) $(SERVER_PKG)

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

docker-up:
	$(COMPOSE) up -d

docker-rebuild:
	$(COMPOSE) up -d --build habit-tracker-service

docker-down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f habit-tracker-service

docker-logs:
	$(COMPOSE) logs -f

db-shell:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

migrate-status:
	docker exec $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "select * from schema_migrations; select to_regclass('public.outbox_event') as outbox_event;"

clean:
	go clean
	-rm -rf $(BIN_DIR)
