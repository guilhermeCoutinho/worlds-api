run-local:
	@PG_URL=postgres://postgres:postgres@localhost:5432/worlds?sslmode=disable go run main.go start

build:
	@docker compose build

up: build
	@docker compose up -d

down:
	@docker compose down

restart:
	@docker compose down
	@docker compose up -d

migrate:
	@docker compose exec app go run migrations/main.go

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@gofmt -w .
	@goimports -w .
