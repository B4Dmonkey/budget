APP_NAME = server
APP_BINARY_DIR = $(PWD)/bin
APP_SCRIPTS_DIR = $(APP_BINARY_DIR)/scripts

DB_DIR=$(PWD)/db
DB_MIGRATION_DIR = $(DB_DIR)/migrations
DB_SCHEMA_FILE = $(DB_DIR)/schema.sql
DB_SCRIPTS_DIR = $(DB_DIR)/scripts
TEST_DB_SCHEMA_FILE = $(DB_DIR)/test_schema.sql

ifneq ("$(wildcard .env)","")
include .env
export $(shell sed 's/=.*//' .env)
endif


# !Not ideal if working multiple people
dev:
	@~/go/bin/air

test: db_migrate_test_database seed run_tests db_clean_up_test_database

run_tests:
	@go test -v ./...

binary:
	@echo "Building $(APP_NAME)..."
	@CGO_ENABLED=0 go build -ldflags="-w -s" -o $(APP_BINARY_DIR)/$(APP_NAME) *.go
	@echo "Build complete."

# Data base migrations
db_build: db_migration_run db_model_sql db_seed_database

db_migration:
	@dbmate --migrations-dir=$(DB_MIGRATION_DIR) new $(name)

db_migration_run:
	@dbmate --migrations-dir=$(DB_MIGRATION_DIR) --schema-file=$(DB_SCHEMA_FILE) up 

db_migration_rollback:
	@dbmate --migrations-dir=$(DB_MIGRATION_DIR) --schema-file=$(DB_SCHEMA_FILE) rollback

db_model_sql sql:
	sqlc generate

db_drop:
	dbmate drop

db_seed_database seed:
	@test -x $(APP_SCRIPTS_DIR)/seed_root_user.sh || chmod +x $(APP_SCRIPTS_DIR)/seed_root_user.sh
	$(APP_SCRIPTS_DIR)/seed_root_user.sh

db_migrate_test_database:
	@dbmate --migrations-dir=$(DB_MIGRATION_DIR) --schema-file=$(TEST_DB_SCHEMA_FILE) --url=$(TEST_DATABASE_URL) up 

db_clean_up_test_database:
	@echo "Cleaning up test database..."
	@rm $(TEST_DB_SCHEMA_FILE) $(DB_DIR)/test.db