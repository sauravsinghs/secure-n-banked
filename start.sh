#!/bin/sh

set -e

if [ -n "$POSTGRES_PASSWORD" ]; then
  export DB_SOURCE="postgresql://root:${POSTGRES_PASSWORD}@postgres:5432/simple_bank?sslmode=disable"
fi

echo "start the app"
exec "$@"
