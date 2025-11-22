#!/bin/sh
set -e

if [ -z "$DB_URL" ]; then
  echo "DB_URL is required for migrations" >&2
  exit 1
fi

echo "Running migrations..."
goose -dir /app/db/migrations postgres "$DB_URL" up

echo "Starting API..."
exec "$@"
