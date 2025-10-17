include .env
export $(shell sed 's/=.*//' .env)

MIGRATIONS_PATH = ./cmd/migrate/migrations
MIGRATE_COMMAND = migrate -path "$(MIGRATIONS_PATH)" -database "$(DB_ADDR)"

.PHONY: create-migration
create-migration:
	@echo "Creating new migration"
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate:
	@echo "Migrating database"
	@$(MIGRATE_COMMAND) up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database"
	@$(MIGRATE_COMMAND) down

.PHONY: migrate-reset
migrate-reset: 
	@echo "Resetting database"
	@$(MIGRATE_COMMAND) drop
	@$(MIGRATE_COMMAND) up

.PHONY: seed-run
seed: 
	@echo "Seeding database"
	@go run cmd/migrate/seed/main.go

.PHONY: seed-reset
seed-reset:
	@echo "Resetting and seeding database"
	@$(MIGRATE_COMMAND) drop
	@$(MIGRATE_COMMAND) up
	@go run cmd/migrate/seed/main.go
