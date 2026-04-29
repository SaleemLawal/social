include .env
MIGRATIONS_PATH= ./cmd/migrate/migrations

.PHONY: migration
migration:
	@echo "Creating migration file..."
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Migrating up..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_URL) up

.PHONY: migrate-down
migrate-down:
	@echo "Migrating down..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_URL) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@echo "Migrating force..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_URL) force $(filter-out $@,$(MAKECMDGOALS))
	@make migrate-up

.PHONY: seed
seed:
	@echo "Seeding database..."
	@DB_URL=${DB_URL} go run cmd/migrate/seed/main.go

.PHONY: run
run: 
	@echo "Starting server..."
	@docker compose up -d --build && docker compose logs -f api

.PHONY: format
format:
	@echo "Formatting code..."
	@gofmt -w .

.PHONY: swagger
swagger:
	@echo "Generating swagger documentation..."
	@swag init -g ./api/main.go -d cmd,internal && swag fmt 