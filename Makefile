include .env
export

DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
DB_PATH="./internal/adapter/storages/postgres/migrations"

migrate-up:
	migrate -path ${DB_PATH} -database ${DB_DSN} up

migrate-up-version:
	@read -p "Enter the version: " version; \
	migrate -path ${DB_PATH} -database ${DB_DSN} up $$version

migrate-down:
	migrate -path ${DB_PATH} -database ${DB_DSN} down

migrate-down-version:
	@read -p "Enter the version: " version; \
	migrate -path ${DB_PATH} -database ${DB_DSN} down $$version

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ${DB_PATH} -seq $$name

migrate-force:
	@read -p "Enter the version: " version; \
	migrate -path ${DB_PATH} -database ${DB_DSN} force $$version

# Run the application
run:
	go run cmd/http/main.go

dev: 
	air

# Start the Docker Compose services
docker-up:
	@docker-compose up -d

# Stop and remove the Docker Compose services
docker-down:
	@docker-compose down

# Clean up Docker resources
docker-clean:
	@docker-compose down -v --remove-orphans

help:
	@echo "  make run         - run application"
	@echo "  make docker-up   - start docker compose"
	@echo "  make docker-down - stop docker compose"
	@echo "  make migrate-up     - run all up migrations"
	@echo "  make migrate-up-version - run up migrations by version"
	@echo "  make migrate-down   - run all down migrations"
	@echo "  make migrate-down-version - run down migrations by version"
	@echo "  make migrate-create - create new migration files"
	@echo "	 make migrate-force " - force migrate version

.PHONY: run docker-up docker-down docker-clean migrate-up migrate-up-version migrate-down migrate-down-version migrate-create migrate-force help
