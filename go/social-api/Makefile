include .env
export $(shell sed 's/=.*//' .env)

MIGRATIONS_PATH = ./cmd/migrate/migrations
MIGRATE_COMMAND = migrate -path "$(MIGRATIONS_PATH)" -database "$(DB_ADDR)"

.PHONY: migrate-create
migrate-create:
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