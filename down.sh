#!/bin/bash

# Option handling
REMOVE_VOLUMES=false
REMOVE_IMAGES=false

# Parse command line arguments
for arg in "$@"
do
    case $arg in
        --volumes|-v)
        REMOVE_VOLUMES=true
        shift
        ;;
        --images|-i)
        REMOVE_IMAGES=true
        shift
        ;;
        --help|-h)
        echo "Usage: ./down.sh [options]"
        echo "Options:"
        echo "  --volumes, -v    Remove volumes (WARNING: This will delete all data)"
        echo "  --images, -i     Remove images"
        echo "  --help, -h       Show this help message"
        exit 0
        ;;
    esac
done

echo "Stopping Hello Pulse development environment..."

if [ "$REMOVE_VOLUMES" = true ] && [ "$REMOVE_IMAGES" = true ]; then
    echo "WARNING: Removing all containers, volumes, and images. This will delete all data."
    docker compose -f docker-compose.dev.yml down --volumes --rmi all
elif [ "$REMOVE_VOLUMES" = true ]; then
    echo "WARNING: Removing all containers and volumes. This will delete all data."
    docker compose -f docker-compose.dev.yml down --volumes
elif [ "$REMOVE_IMAGES" = true ]; then
    echo "Removing all containers and images."
    docker compose -f docker-compose.dev.yml down --rmi all
else
    echo "Stopping containers (preserving data)."
    docker compose -f docker-compose.dev.yml down
fi

echo "Environment is down."