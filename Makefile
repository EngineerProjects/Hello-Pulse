# Hello Pulse - Simplified Makefile

# Variables
SHELL := /bin/bash
DOCKER_COMPOSE_DEV := docker compose -f docker-compose.dev.yml
DOCKER_COMPOSE_PROD := docker compose -f docker-compose.yml

# Colors for terminal output
COLOR_RESET := \033[0m
COLOR_GREEN := \033[0;32m
COLOR_YELLOW := \033[0;33m
COLOR_CYAN := \033[0;36m
COLOR_RED := \033[0;31m

# Phony targets
.PHONY: help setup dev prod down logs clean \
        prune-all prune-containers prune-images prune-volumes prune-networks \
        prune-builders backend-test db-setup gen-jwt

# Default target
help:
	@echo -e "$(COLOR_CYAN)Hello Pulse - Available Commands:$(COLOR_RESET)"
	@echo -e "$(COLOR_GREEN)Main Commands:$(COLOR_RESET)"
	@echo "  make setup          - Initial project setup (create dirs, .env file)"
	@echo "  make gen-jwt        - Generate a secure JWT secret and update .env"
	@echo "  make dev            - Start development environment"
	@echo "  make prod           - Start production environment"
	@echo "  make down           - Stop all containers"
	@echo "  make logs           - View logs from all services"
	@echo -e "$(COLOR_GREEN)Cleanup Commands:$(COLOR_RESET)"
	@echo "  make clean          - Stop containers and remove volumes (preserves images)"
	@echo "  make prune-all      - Prune everything (containers, images, volumes, networks, builders)"
	@echo "  make prune-containers - Prune only containers"
	@echo "  make prune-images   - Prune only images"
	@echo "  make prune-volumes  - Prune only volumes"
	@echo "  make prune-networks - Prune only networks"
	@echo "  make prune-builders - Prune builder cache"
	@echo -e "$(COLOR_GREEN)Testing:$(COLOR_RESET)"
	@echo "  make backend-test   - Run backend tests"
	@echo -e "$(COLOR_GREEN)Database:$(COLOR_RESET)"
	@echo "  make db-setup       - Initialize database extensions"
	@echo -e "$(COLOR_GREEN)Linting:$(COLOR_RESET)"
	@echo "  make check-lint     - Check code linting"

# Setup command
setup:
	@echo -e "$(COLOR_CYAN)Setting up Hello Pulse project...$(COLOR_RESET)"
	@if [ ! -f .env ]; then \
		echo -e "$(COLOR_YELLOW)Creating .env file from .env.example...$(COLOR_RESET)"; \
		cp .env.example .env; \
		echo -e "$(COLOR_YELLOW)Please review and configure the .env file.$(COLOR_RESET)"; \
	fi
	@mkdir -p volumes/postgres-data volumes/pgadmin-data volumes/minio-data
	@chmod +x up.sh down.sh scripts/check-lint.sh
	@echo -e "$(COLOR_GREEN)Setup complete!$(COLOR_RESET)"

# Generate JWT Secret
gen-jwt:
	@echo -e "$(COLOR_CYAN)Generating secure JWT secret...$(COLOR_RESET)"
	@# Check if OpenSSL is installed
	@if ! command -v openssl &> /dev/null; then \
		echo -e "$(COLOR_RED)OpenSSL not found. Please install it manually.$(COLOR_RESET)"; \
		exit 1; \
	fi
	@# Generate the JWT secret
	@SECRET=$$(openssl rand -base64 64 | tr -d '\n') && \
	if [ ! -f .env ]; then \
		cp .env.example .env; \
	fi && \
	NEW_CONTENT=$$(grep -v "^JWT_SECRET=" .env) && \
	echo "$$NEW_CONTENT" > .env && \
	echo "JWT_SECRET=$$SECRET" >> .env && \
	echo -e "$(COLOR_GREEN)JWT_SECRET has been updated in your .env file.$(COLOR_RESET)"

# Environment commands
dev:
	@echo -e "$(COLOR_CYAN)Starting development environment...$(COLOR_RESET)"
	@./up.sh

prod:
	@echo -e "$(COLOR_CYAN)Starting production environment...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE_PROD) up -d

down:
	@echo -e "$(COLOR_CYAN)Stopping all services...$(COLOR_RESET)"
	@./down.sh

clean:
	@echo -e "$(COLOR_RED)WARNING: This will remove all containers and volumes.$(COLOR_RESET)"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "$(COLOR_YELLOW)Removing all containers and volumes...$(COLOR_RESET)"; \
		$(DOCKER_COMPOSE_DEV) down -v; \
		echo -e "$(COLOR_GREEN)Clean complete.$(COLOR_RESET)"; \
	else \
		echo -e "$(COLOR_YELLOW)Clean operation cancelled.$(COLOR_RESET)"; \
	fi

logs:
	@echo -e "$(COLOR_CYAN)Viewing logs from all services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE_DEV) logs -f

# Docker pruning commands
prune-all:
	@echo -e "$(COLOR_RED)WARNING: This will prune everything (containers, images, volumes, networks, builders).$(COLOR_RESET)"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "$(COLOR_YELLOW)Pruning all Docker resources...$(COLOR_RESET)"; \
		docker system prune -af --volumes; \
		echo -e "$(COLOR_GREEN)Prune complete.$(COLOR_RESET)"; \
	else \
		echo -e "$(COLOR_YELLOW)Prune operation cancelled.$(COLOR_RESET)"; \
	fi

prune-containers:
	@echo -e "$(COLOR_CYAN)Pruning containers...$(COLOR_RESET)"
	@docker container prune -f
	@echo -e "$(COLOR_GREEN)Container prune complete.$(COLOR_RESET)"

prune-images:
	@echo -e "$(COLOR_CYAN)Pruning images...$(COLOR_RESET)"
	@docker image prune -af
	@echo -e "$(COLOR_GREEN)Image prune complete.$(COLOR_RESET)"

prune-volumes:
	@echo -e "$(COLOR_RED)WARNING: This will remove all unused volumes and their data.$(COLOR_RESET)"
	@read -p "Are you sure you want to continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "$(COLOR_YELLOW)Pruning volumes...$(COLOR_RESET)"; \
		docker volume prune -f; \
		echo -e "$(COLOR_GREEN)Volume prune complete.$(COLOR_RESET)"; \
	else \
		echo -e "$(COLOR_YELLOW)Volume prune cancelled.$(COLOR_RESET)"; \
	fi

prune-networks:
	@echo -e "$(COLOR_CYAN)Pruning networks...$(COLOR_RESET)"
	@docker network prune -f
	@echo -e "$(COLOR_GREEN)Network prune complete.$(COLOR_RESET)"

prune-builders:
	@echo -e "$(COLOR_CYAN)Pruning builder cache...$(COLOR_RESET)"
	@docker builder prune -af
	@echo -e "$(COLOR_GREEN)Builder cache prune complete.$(COLOR_RESET)"

# Testing commands
backend-test:
	@echo -e "$(COLOR_CYAN)Running backend tests...$(COLOR_RESET)"
	@cd backend && go test -v ./...

# Database commands
db-setup:
	@echo -e "$(COLOR_CYAN)Setting up database extensions...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE_DEV) exec db psql -U $${POSTGRES_USER:-postgres} -d $${POSTGRES_DB_NAME:-hellopulsedb} -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
	@$(DOCKER_COMPOSE_DEV) exec db psql -U $${POSTGRES_USER:-postgres} -d $${POSTGRES_DB_NAME:-hellopulsedb} -c "CREATE EXTENSION IF NOT EXISTS vector;"
	@echo -e "$(COLOR_GREEN)Database setup complete.$(COLOR_RESET)"

# Check lint command
check-lint:
	@echo -e "$(COLOR_CYAN)Checking code lint...$(COLOR_RESET)"
	@./scripts/check-lint.sh
	@echo -e "$(COLOR_GREEN)Lint check complete.$(COLOR_RESET)"