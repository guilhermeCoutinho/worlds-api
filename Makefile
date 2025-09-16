setup:
	@go mod download && go mod tidy

run-local:
	@PG_URL=postgres://postgres:postgres@localhost:5432/worlds?sslmode=disable \
	REDIS_URL=redis://localhost:6379 go run main.go start --verbose=5

build:
	@docker compose build

up: build
	@docker compose up -d

down:
	@docker compose down

restart:
	@docker compose down
	@docker compose up -d

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@gofmt -w .
	@goimports -w .

test: migrate-reset migrate
	@go test ./test/end2end/...


## migrate: execute the postgres migration; use ADDRESS="" and PASSWORD="" to specify another database
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
