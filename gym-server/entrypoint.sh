#!/bin/sh

echo "Running migrations..."
./migrate \
  --storage-url="$DB_URL" \
  --migrations-path="./migrations"

echo "Starting server..."
exec ./main --config=./config/dev.yaml