.PHONY: help install dev-backend dev-frontend sync build-frontend build-backend build clean swagger jet

# Load environment variables
include .env
export

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	cd pkg && go mod download
	@echo "Installing frontend dependencies..."
	cd web && bun install
	@echo "Dependencies installed!"

dev-backend: ## Start backend development server
	@echo "Starting backend on :8088..."
	cd pkg && go run main.go

dev-frontend: ## Start frontend development server
	@echo "Starting frontend on :3000..."
	cd web && bun dev

sync: ## Sync videos from YouTube (requires YOUTUBE_API_KEY)
	@bash scripts/sync-youtube.sh

build-frontend: ## Build frontend for production
	@echo "Building frontend..."
	cd web && bun run build
	@echo "Frontend built to web/.output/public/"

build-backend: ## Build backend binary
	@echo "Building backend..."
	cd pkg && go build -o zeedzad
	@echo "Backend binary created: pkg/zeedzad"

build: build-frontend build-backend ## Build both frontend and backend
	@echo "Copying frontend assets to embed location..."
	@mkdir -p pkg/web/public
	@cp -r web/.output/public/* pkg/web/public/
	@echo "Build complete!"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf web/.output
	@rm -rf web/.nuxt
	@rm -rf pkg/zeedzad
	@rm -rf pkg/web/public/*
	@echo "Clean complete!"

swagger: ## Generate Swagger API documentation
	@echo "Generating Swagger docs..."
	cd pkg && swag init -d . -g server/server.go -o docs --ot go
	@echo "Swagger docs generated!"

jet: ## Generate Go-Jet database models
	@echo "Generating Go-Jet models..."
	cd pkg && jet -dsn=$(SQLITE_PATH) -schema=file://db/schema.sql -path=./repository/table
	@echo "Go-Jet models generated!"

db-init: ## Initialize database with schema
	@echo "Initializing database..."
	@mkdir -p data
	@sqlite3 $(SQLITE_PATH) < pkg/db/schema.sql
	@echo "Database initialized at $(SQLITE_PATH)"

db-reset: ## Reset database (WARNING: deletes all data)
	@echo "Resetting database..."
	@rm -f $(SQLITE_PATH)
	@$(MAKE) db-init
	@echo "Database reset complete!"

run: ## Run production binary
	@cd pkg && ./zeedzad

# Development workflow shortcuts
.PHONY: start setup

setup: install db-init ## Complete setup for new developers
	@echo ""
	@echo "Setup complete! Next steps:"
	@echo "1. Copy .env.example to .env and configure"
	@echo "2. Add your YOUTUBE_API_KEY to .env"
	@echo "3. Run 'make sync' to fetch videos"
	@echo "4. Run 'make dev-backend' in one terminal"
	@echo "5. Run 'make dev-frontend' in another terminal"
	@echo "6. Open http://localhost:3000 in your browser"

start: ## Start both backend and frontend (requires tmux)
	@echo "Starting services with tmux..."
	@tmux new-session -d -s zeedzad 'cd pkg && go run main.go'
	@tmux split-window -h 'cd web && bun dev'
	@tmux attach -t zeedzad
