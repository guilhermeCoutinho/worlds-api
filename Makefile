.PHONY: test setup run-local up restart migrate-init migrate migrate-reset

setup:
	@go mod download && go mod tidy

run-local:
	@PG_URL=postgres://postgres:postgres@localhost:5432/worlds?sslmode=disable \
	REDIS_URL=redis://localhost:6379 go run main.go start --verbose=5

up: 
	@docker compose up -d --build

restart:
	@docker compose down
	@docker compose up -d --build

test:
	@echo "Recreating database for tests"
	@make migrate-reset > /dev/null
	@make migrate > /dev/null
	@go test ./test/end2end/... -count=1 -v


ADDRESS := localhost:5432
PASSWORD := "postgres"
USER := postgres
migrate-init:
	@echo "Running migration on database: $(ADDRESS)"
	@cd migrations && go run *.go init -address $(ADDRESS) -pass $(PASSWORD) -user $(USER)

migrate:
	@echo "Running migration on database: $(ADDRESS)"
	@cd migrations && go run *.go -address $(ADDRESS) -pass $(PASSWORD) -user $(USER)

migrate-reset:
	@echo "Running migration on database: $(ADDRESS)"
	@cd migrations && go run *.go -address $(ADDRESS) -pass $(PASSWORD) -user $(USER) reset
