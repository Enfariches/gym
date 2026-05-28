#!/bin/sh

echo "Building binaries..."
cd /gym || exit 1

go mod tidy

go build -o migrate ./cmd/migrator/migrator.go || {
  echo "Failed to build migrate"
  exit 1
}

go build -o main ./cmd/gymserver/main.go || {
  echo "Failed to build main"
  exit 1
}

echo "Running migrations..."
./migrate \
  --storage-url="$DB_URL" \
  --migrations-path="./migrations"

echo "Starting server..."
exec ./main --config=./config/dev.yaml
