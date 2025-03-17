#!/bin/bash

# hello-pulse project setup script
# This script creates the entire project structure with appropriate files

set -e  # Exit on error

# Create root directory
mkdir -p hello-pulse
cd hello-pulse

# Create main directories
mkdir -p cmd internal pkg ai-services clients infrastructure scripts docs .github

# Create cmd subdirectories and Go files
cmd_dirs=("api-gateway" "auth-service" "user-service" "project-service" "file-service" "collaboration-service" "ai-orchestrator" "db-migrator")
for dir in "${cmd_dirs[@]}"; do
    mkdir -p "cmd/$dir"
    echo "package main

func main() {
    // $dir entry point
}
" > "cmd/$dir/main.go"
done

# Create internal directory structure and Go files
mkdir -p internal/models internal/database/migrations internal/config internal/api/middleware internal/api/response internal/api/validation
mkdir -p internal/gateway/handlers internal/gateway/proxy internal/gateway/routes
mkdir -p internal/auth/handlers internal/auth/repository internal/auth/services
mkdir -p internal/user/handlers internal/user/repository internal/user/services
mkdir -p internal/project/handlers internal/project/repository internal/project/services
mkdir -p internal/file/handlers internal/file/repository internal/file/services
mkdir -p internal/collaboration/handlers internal/collaboration/repository internal/collaboration/services
mkdir -p internal/ai/handlers internal/ai/repository internal/ai/clients internal/ai/services

# Create model Go files
model_files=("user.go" "token.go" "organization.go" "project.go" "file.go" "conversation.go" "message.go" "embedding.go")
for file in "${model_files[@]}"; do
    echo "package models

// $(basename "$file" .go) model
" > "internal/models/$file"
done

# Create database files
echo "package database

// Connection setup
" > "internal/database/connection.go"

echo "package database

// Repository pattern
" > "internal/database/repository.go"

# Create database migration files
migration_files=("001_initial_schema.go" "002_users_auth.go" "003_projects.go" "004_files.go" "005_ai_components.go")
for file in "${migration_files[@]}"; do
    echo "package migrations

// $(basename "$file" .go) migration
" > "internal/database/migrations/$file"
done

# Create config files
echo "package config

// Configuration loading
" > "internal/config/config.go"

echo "package config

// Configuration validation
" > "internal/config/validation.go"

# Create API middleware files
middleware_files=("auth.go" "cors.go" "logging.go")
for file in "${middleware_files[@]}"; do
    echo "package middleware

// $(basename "$file" .go) middleware
" > "internal/api/middleware/$file"
done

# Create API response files
echo "package response

// Error response handling
" > "internal/api/response/error.go"

echo "package response

// Success response handling
" > "internal/api/response/success.go"

# Create API validation file
echo "package validation

// Request validator
" > "internal/api/validation/validator.go"

# Create gateway files
echo "package routes

// Route setup
" > "internal/gateway/routes/routes.go"

# Create proxy files
proxy_files=("auth_proxy.go" "user_proxy.go" "ws_proxy.go")
for file in "${proxy_files[@]}"; do
    echo "package proxy

// $(basename "$file" .go) proxy
" > "internal/gateway/proxy/$file"
done

# Create auth service files
echo "package handlers

// Authentication endpoints
" > "internal/auth/handlers/auth_handler.go"

echo "package handlers

// Auth routes setup
" > "internal/auth/handlers/routes.go"

echo "package repository

// User repository
" > "internal/auth/repository/user_repo.go"

echo "package repository

// Token repository
" > "internal/auth/repository/token_repo.go"

echo "package services

// Authentication logic
" > "internal/auth/services/auth_service.go"

echo "package services

// JWT handling
" > "internal/auth/services/jwt_service.go"

# Create AI files
echo "package handlers

// AI endpoints
" > "internal/ai/handlers/ai_handler.go"

echo "package repository

// Conversation repository
" > "internal/ai/repository/conversation_repo.go"

echo "package repository

// Message repository
" > "internal/ai/repository/message_repo.go"

# Create AI client files
ai_client_files=("nlp_client.go" "rag_client.go" "vision_client.go")
for file in "${ai_client_files[@]}"; do
    echo "package clients

// $(basename "$file" .go) client
" > "internal/ai/clients/$file"
done

echo "package services

// AI orchestration logic
" > "internal/ai/services/orchestrator.go"

# Create file service storage file
echo "package services

// Storage handling
" > "internal/file/services/storage_service.go"

# Create collaboration websocket file
echo "package services

// WebSocket handling
" > "internal/collaboration/services/ws_service.go"

# Create public packages
mkdir -p pkg/jwt pkg/httputils pkg/logger pkg/vectordb

echo "package jwt

// JWT helper functions
" > "pkg/jwt/jwt.go"

echo "package httputils

// Common HTTP utilities
" > "pkg/httputils/http.go"

echo "package logger

// Logger setup
" > "pkg/logger/logger.go"

echo "package vectordb

// Vector DB client
" > "pkg/vectordb/client.go"

# Create AI services directories
ai_services=("nlp-service" "rag-service" "vision-service" "speech-service" "web-search" "image-generation")

# Create the single Dockerfile for Python services with minimal content
mkdir -p ai-services
cat > "ai-services/Dockerfile" << 'EOF'
FROM python:3.9-slim
WORKDIR /app
CMD ["python", "-m", "app.main"]
EOF

# Create Python service structure and __init__.py files
for service in "${ai_services[@]}"; do
    mkdir -p "ai-services/$service/app/api/routes" "ai-services/$service/app/core" "ai-services/$service/app/services"
    
    # Create __init__.py files
    touch "ai-services/$service/app/__init__.py"
    touch "ai-services/$service/app/api/__init__.py"
    touch "ai-services/$service/app/api/routes/__init__.py"
    touch "ai-services/$service/app/core/__init__.py"
    touch "ai-services/$service/app/services/__init__.py"
    
    # Create empty main.py for each service
    touch "ai-services/$service/app/main.py"
    
    # Create empty requirements.txt
    touch "ai-services/$service/requirements.txt"
done

# Create clients directories
mkdir -p clients/web/public clients/web/src/{components,pages,services,hooks,contexts,utils}
mkdir -p clients/mobile/src clients/desktop/src

# Create empty package.json files for clients
for client in "web" "mobile" "desktop"; do
    touch "clients/$client/package.json"
done

# Create empty web client essential files
touch "clients/web/src/services/api.ts"
touch "clients/web/src/services/auth.ts"
touch "clients/web/src/services/websocket.ts"

# Create minimal web Dockerfile
cat > "clients/web/Dockerfile" << EOF
FROM node:16-alpine
WORKDIR /app
CMD ["npm", "start"]
EOF

# Create infrastructure directories
mkdir -p infrastructure/{docker/{development,production},postgres/init-scripts,vector-db,redis,kubernetes/{base,overlays/{development,production}}}

# Create empty docker-compose files
touch "infrastructure/docker/development/docker-compose.yml"
touch "infrastructure/docker/production/docker-compose.yml"

# Create empty database and config files
touch "infrastructure/postgres/init-scripts/init-db.sh"
chmod +x "infrastructure/postgres/init-scripts/init-db.sh"
touch "infrastructure/vector-db/config.yaml"
touch "infrastructure/redis/redis.conf"

# Create empty Kubernetes manifest files
touch "infrastructure/kubernetes/base/api-gateway.yaml"
touch "infrastructure/kubernetes/base/auth-service.yaml"

# Create minimal script files
touch "scripts/setup.sh"
chmod +x "scripts/setup.sh"

touch "scripts/dev.sh"
chmod +x "scripts/dev.sh"

touch "scripts/migrate.sh"
chmod +x "scripts/migrate.sh"

touch "scripts/test.sh"
chmod +x "scripts/test.sh"

touch "scripts/deploy.sh"
chmod +x "scripts/deploy.sh"

# Create empty documentation files
mkdir -p "docs/api" "docs/diagrams" "docs/development"
touch "docs/architecture.md"
touch "docs/api/openapi.yaml"
touch "docs/development/getting-started.md"

# Create empty GitHub workflow files
mkdir -p .github/workflows
touch ".github/workflows/ci.yml"
touch ".github/workflows/deploy.yml"

# Create empty root files
cat > "go.mod" << EOF
module hello-pulse

go 1.19
EOF

touch "Makefile"
touch ".env.example"
touch "docker-compose.yml"
touch "README.md"

echo "Project structure created successfully!"