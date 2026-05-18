# Load .env variables
include .env
export

# Build DB URL dynamically
DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}

APP_NAME=qonevo-backend
MAIN_PATH=cmd/server/main.go

# -----------------------
# 🚀 APP COMMANDS
# -----------------------

run:
	go run $(MAIN_PATH)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

tidy:
	go mod tidy

# -----------------------
# 🗄️ MIGRATION COMMANDS
# -----------------------

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-reset:
	migrate -path migrations -database "$(DB_URL)" down

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

migrate-force:
	@read -p "Enter version: " version; \
	migrate -path migrations -database "$(DB_URL)" force $$version