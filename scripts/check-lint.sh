#!/bin/sh

set -o xtrace

cd backend && golangci-lint run \
    --enable gocritic \
    --enable goimports \
    --enable exhaustive \
    --enable gosec \
    --enable goimports \
    --enable nilerr \
    --enable gofumpt \
    --enable revive \
    ./...

cd ..

cd clients/web && npm run lint
