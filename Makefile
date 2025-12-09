APP_NAME=horizonx-server
ENTRY=./cmd/horizonx-server
MIGRATE_SRC=cmd/migrate/main.go
SEED_SRC=cmd/seed/main.go

BIN_DIR=bin
MIGRATE_BIN=$(BIN_DIR)/$(APP_NAME)-migrate
SEED_BIN=$(BIN_DIR)/$(APP_NAME)-seed

DB_PATH=horizonx.db
MIGRATION_DIR=internal/repository/sqlite/migrations

# --- App Commands ---
build:
	@echo "Building binaries..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(ENTRY)
	@go build -o $(MIGRATE_BIN) $(MIGRATE_SRC)
	@go build -o $(SEED_BIN) $(SEED_SRC)
	@echo "Build complete! Check $(BIN_DIR)/"

run:
	@go run $(ENTRY)

clean:
	@rm -rf $(BIN_DIR)
	@rm -f $(DB_PATH)

# --- Database Commands (Development) ---
migrate-up:
	@go run $(MIGRATE_SRC) -op=up -db=$(DB_PATH)

migrate-down:
	@go run $(MIGRATE_SRC) -op=down -steps=1 -db=$(DB_PATH)

migrate-fresh:
	@echo "Resetting database..."
	@go run $(MIGRATE_SRC) -op=down -db=$(DB_PATH)
	@go run $(MIGRATE_SRC) -op=up -db=$(DB_PATH)
	@echo "Database fresh and clean!"

migrate-version:
	@go run $(MIGRATE_SRC) -op=version -db=$(DB_PATH)

migrate-create:
	@test -n "$(name)" || (echo "Error: name is required"; exit 1)
	@echo "Creating migration files..."
	@mkdir -p $(MIGRATION_DIR)
	@touch $(MIGRATION_DIR)/$$(date +%Y%m%d%H%M%S)_$(name).up.sql
	@touch $(MIGRATION_DIR)/$$(date +%Y%m%d%H%M%S)_$(name).down.sql
	@echo "Files created in $(MIGRATION_DIR)"

seed:
	@go run $(SEED_SRC) -db=$(DB_PATH)

.PHONY: build run clean migrate-up migrate-down migrate-fresh migrate-version migrate-create seed
