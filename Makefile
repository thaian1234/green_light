include .env
export

DB_DSN="postgres://${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
DB_PATH="./internal/adapter/storages/postgres/migrations"

migrate-up:
	migrate -path "${DB_PATH}" -database "${DB_DSN}" up

migrate-down:
	migrate -path "${DB_PATH}" -database "${DB_DSN}" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir "${DB_PATH}" -seq $$name

# Run the application
run:
	go run cmd/http/main.go

# Start the Docker Compose services
docker-up:
	@docker-compose up -d

# Stop and remove the Docker Compose services
docker-down:
	@docker-compose down

help:
	@echo "  make run         - run application"
	@echo "  make docker-up   - start docker compose"
	@echo "  make docker-down - stop docker compose"
	@echo "  make migrate-up     - run all up migrations"
	@echo "  make migrate-down   - run all down migrations"
	@echo "  make migrate-create - create new migration files"
