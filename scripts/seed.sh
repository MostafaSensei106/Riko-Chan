#!/bin/bash

# Database seeding script
set -e


DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-future_messages}

echo "Seeding database with test data..."

# Create test user
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << EOF
INSERT INTO users (id, first_name, language, timezone, created_at, updated_at)
VALUES (123456789, 'Test User', 'en', 'UTC', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
EOF

echo "Database seeding completed!"
