# =========================================
# Load Environment Variables
# =========================================

include .env
export

# =========================================
# App Configuration
# =========================================

APP_NAME=qonevo-backend
MAIN_PATH=./cmd/server/main.go
BUILD_DIR=./bin

# =========================================
# Database Configuration
# =========================================

DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}

# =========================================
# Go Commands
# =========================================

GO=go

# =========================================
# App Commands
# =========================================

.PHONY: run dev build clean tidy fmt vet

run:
	$(GO) run $(MAIN_PATH)

dev:
	air

build:
	mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

clean:
	rm -rf $(BUILD_DIR)
	rm -rf tmp

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

# =========================================
# Migration Commands
# =========================================

.PHONY: migrate-up migrate-down migrate-reset migrate-version migrate-force migrate-create

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-reset:
	migrate -path migrations -database "$(DB_URL)" down -all

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

migrate-force:
	@read -p "Enter migration version: " version; \
	migrate -path migrations -database "$(DB_URL)" force $$version

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# =========================================
# Database Utilities
# =========================================

.PHONY: db-url

db-url:
	@echo $(DB_URL)

# =========================================
# Production Commands
# =========================================

.PHONY: prod

prod:
	$(BUILD_DIR)/$(APP_NAME)