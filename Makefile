# --- Variables ---
APP_NAME=horizonx-server
AGENT_NAME=horizonx-agent

# Entry Points
SERVER_ENTRY=./cmd/server/main.go
AGENT_ENTRY=./cmd/agent/main.go
MIGRATE_SRC=./cmd/migrate/main.go
SEED_SRC=./cmd/seed/main.go

# Binaries Output
BIN_DIR=bin
SERVER_BIN=$(BIN_DIR)/server
AGENT_BIN=$(BIN_DIR)/agent
MIGRATE_BIN=$(BIN_DIR)/migrate
SEED_BIN=$(BIN_DIR)/seed

# Configs
MIGRATION_DIR=internal/storage/postgres/migrations

# --- Build Commands ---
build:
	@echo "🚧 Building binaries..."
	@mkdir -p $(BIN_DIR)
	@echo "   • Compiling Server..."
	@go build -o $(SERVER_BIN) $(SERVER_ENTRY)
	@echo "   • Compiling Agent..."
	@go build -o $(AGENT_BIN) $(AGENT_ENTRY)
	@echo "   • Compiling Tools (Migrate & Seed)..."
	@go build -o $(MIGRATE_BIN) $(MIGRATE_SRC)
	@go build -o $(SEED_BIN) $(SEED_SRC)
	@echo "✅ Build complete! Check $(BIN_DIR)/"

# Run Server (Dev Mode)
run-server:
	@go run $(SERVER_ENTRY)

# Run Agent (Dev Mode - might need sudo for hardware stats)
run-agent:
	@go run $(AGENT_ENTRY)

clean:
	@rm -rf $(BIN_DIR)
	@echo "🧹 Binaries cleaned."

# --- Database Commands (Postgres) ---
migrate-up:
	@go run $(MIGRATE_SRC) -op=up

migrate-down:
	@go run $(MIGRATE_SRC) -op=down -steps=1

migrate-fresh:
	@echo "🧨 Resetting database..."
	@go run $(MIGRATE_SRC) -op=down # Revert all
	@go run $(MIGRATE_SRC) -op=up   # Apply all
	@echo "✨ Database fresh and clean!"

migrate-version:
	@go run $(MIGRATE_SRC) -op=version

# Example: make migrate-create name=init_schema
migrate-create:
	@test -n "$(name)" || (echo "Error: name is required. Usage: make migrate-create name=something"; exit 1)
	@echo "Creating migration files..."
	@mkdir -p $(MIGRATION_DIR)
	@touch $(MIGRATION_DIR)/$$(date +%Y%m%d%H%M%S)_$(name).up.sql
	@touch $(MIGRATION_DIR)/$$(date +%Y%m%d%H%M%S)_$(name).down.sql
	@echo "📄 Files created in $(MIGRATION_DIR)"

seed:
	@go run $(SEED_SRC)

.PHONY: build run-server run-agent clean migrate-up migrate-down migrate-fresh migrate-version migrate-create seed
