# Configuration
# ==================================
DB_DRIVER ?= postgresql
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432
DB_NAME ?= db_apollo
DB_USER ?= posgres_user
DB_PASSWORD ?= postgres123
DB_SSL ?= disable

MIGRATION_EXT ?= sql
MIGRATION_PATH ?= migrations/
MIGRATION_DSN ?= $(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

SWAGGER_OUTPUT ?= ./docs
TEST_PACKAGES ?= ./...

# Commands
# ==============================================
.PHONY: run test migrate-create auto-migrate generate-api-doc help

help:  ## Display this help message
	@echo.
	@echo Usage:
	@echo   make ^<target^>
	@echo.
	@echo Targets:
	@echo   run               Run the HTTP server
	@echo   test              Run tests with race detection
	@echo   test-coverage     Run tests and show coverage report
	@echo   migrate-create    Create new migration file (NAME required)
	@echo   auto-migrate      Run all pending migrations
	@echo   migrate-up        Apply all up migrations
	@echo   migrate-down      Rollback the last migration
	@echo   migrate-status    Show migration status
	@echo   generate-api-doc  Generate API documentation
	@echo   build            Build the application
	@echo   clean            Clean build artifacts
	@echo   help             Display this help message

##@ Development

run: ## Run the HTTP server
	@echo "Starting HTTP server..."
	@air

test: ## Run tests with race detection
	@echo "Running tests..."
	@go test -v -race $(TEST_PACKAGES) -coverprofile=coverage.out
	@go tool cover -func=coverage.out

test-coverage: test ## Run tests and show coverage report
	@go tool cover -html=coverage.out -o coverage.html

##@ Database

migrate-create: ## Create new migration file (NAME required)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME parameter is required"; \
		echo "Usage: make migrate-create NAME=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating migration '$(NAME)'..."
	@migrate create -ext $(MIGRATION_EXT) -dir $(MIGRATION_PATH) $(NAME)

auto-migrate: ## Run all pending migrations
	@echo "Running migrations..."
	@go run cmd/migrations.go

migrate-up: ## Apply all up migrations
	@migrate -path $(MIGRATION_PATH) -database "$(MIGRATION_DSN)" up

migrate-down: ## Rollback the last migration
	@migrate -path $(MIGRATION_PATH) -database "$(MIGRATION_DSN)" down 1

migrate-status: ## Show migration status
	@migrate -path $(MIGRATION_PATH) -database "$(MIGRATION_DSN)" version

##@ Documentation

generate-api-doc: ## Generate API documentation
	@echo "Generating API documentation..."
	@swag init \
		--generalInfo ./cmd/http/http.go \
		--parseInternal \
		--parseDepth 1 \
		--propertyStrategy pascalcase \
		--outputTypes go,json,yaml \
		--output $(SWAGGER_OUTPUT)
	@swag fmt

run-godoc:
	@godoc -http=:6060

##@ Build

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/server cmd/http/http.go

clean: ## Clean build artifacts
	@rm -rf bin/ coverage.*

##@ Generate Dependency Injection
wire:
	@wire ./modules/auth \
	./modules/user \
	./modules/country
