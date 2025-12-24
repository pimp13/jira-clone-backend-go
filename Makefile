MAIN_FILE = ./cmd/server/main.go
BUILD_FILE = ./bin/server.exe
SWAGGER_DOCS_DIR = ./docs
WIRE_LOCATION = ./internal/infrastructure/di

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m # No Color

.PHONY: run help dev swag build swag-fmt clean gofmt test migrate

run: build ## Run server in production mode
	@$(BUILD_FILE)

help: ## Show help
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(GREEN) %s\n", $$1, $$2}'

dev: ## Run server in development mode
	@echo "Running server in development mode"
	@air -c .air.toml

swag: swag-fmt ## Create swagger document
	@echo "Creating swagger document..."
	@swag init -g $(MAIN_FILE) --output $(SWAGGER_DOCS_DIR)

build: swag ## Build application binary and swagger docs (includes swag-fmt)
	@echo "Building application..."
	@go build -o $(BUILD_FILE) $(MAIN_FILE)

swag-fmt: gofmt ## Format swagger docs
	@echo "Formatting swagger docs"
	@swag fmt

clean: ## Clean up docs and tmp files
	@rm -rf ./docs
	@rm -rf ./tmp

gofmt: ## Format go files
	@echo "Formatting go files"
	@go fmt ./...

test: ## Run tests
	@echo "Running tests"
	@go test ./...

wire: ## Generate dependency injection code
	@echo "Generating dependency injection code..."
	@cd $(WIRE_LOCATION) && wire
	@echo "Dependency injection code generated."


me:
	@echo " ____                         ";
	@echo "|  _ \ ___  _   _ _   _  __ _ ";
	@echo "| |_) / _ \| | | | | | |/ _` |";
	@echo "|  __/ (_) | |_| | |_| | (_| |";
	@echo "|_|   \___/ \__,_|\__, |\__,_|";
	@echo "                  |___/       ";




migrate:
	@go run entgo.io/ent/cmd/ent generate ./ent/schema