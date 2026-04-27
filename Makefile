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