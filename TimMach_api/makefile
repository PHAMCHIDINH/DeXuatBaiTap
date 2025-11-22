# ---------------------------
# CONFIG
# ---------------------------

-include .env

DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= heartdb
DB_HOST ?= localhost
DB_PORT ?= 5432

GOOSE_DIR ?= ./db/migrations

DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# ---------------------------
# GOOSE COMMANDS
# ---------------------------

.PHONY: goose-up
goose-up:
	goose -dir $(GOOSE_DIR) postgres "$(DB_URL)" up

.PHONY: goose-down
goose-down:
	goose -dir $(GOOSE_DIR) postgres "$(DB_URL)" down

.PHONY: goose-status
goose-status:
	goose -dir $(GOOSE_DIR) postgres "$(DB_URL)" status

.PHONY: goose-create
goose-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Missing migration name. Usage: make goose-create name=add_users_table"; \
		exit 1; \
	fi
	goose -dir $(GOOSE_DIR) create $(name) sql

# ---------------------------
# SQLC COMMAND
# ---------------------------

.PHONY: sqlc
sqlc:
	sqlc generate

# ---------------------------
# HELP
# ---------------------------

.PHONY: help
help:
	@echo ""
	@echo "📌 USEFUL COMMANDS:"
	@echo ""
	@echo "  make goose-up             - Apply all migrations"
	@echo "  make goose-down           - Roll back last migration"
	@echo "  make goose-status         - Check migration status"
	@echo "  make goose-create name=X  - Create new migration file"
	@echo "  make sqlc                 - Generate Go code from SQL queries"
	@echo ""
