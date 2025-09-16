run-local:
	@go run main.go start

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