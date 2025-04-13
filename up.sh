#!/bin/bash

# Check if .env file exists, if not, create from example
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please review and update the .env file with your configuration."
    exit 1
fi

# Build and start the development environment
echo "Starting Hello Pulse development environment..."
docker compose -f docker-compose.dev.yml up -d --build

echo "Services:"
echo "- Frontend: http://localhost:3000"
echo "- Backend API: http://localhost:5000"
echo "- API via Nginx: http://localhost:7000"
echo "- MinIO Console: http://localhost:8900"
echo "- PGAdmin: http://localhost:8888"
echo "  (login with credentials from .env file)"