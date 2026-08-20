#!/bin/sh
set -e

# This script is run by PostgreSQL on container startup
# It creates databases from the DB_NAMES environment variable

if [ -z "$DB_NAMES" ]; then
    echo "No DB_NAMES environment variable set, skipping database creation"
    exit 0
fi

echo "Creating databases from DB_NAMES: $DB_NAMES"

# Split DB_NAMES by comma — POSIX-portable (no bash arrays / here-strings)
OLD_IFS=$IFS
IFS=','
for db in $DB_NAMES; do
    db=$(echo "$db" | xargs)
    if [ -n "$db" ]; then
        echo "Creating database: $db"
        psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
            SELECT 'CREATE DATABASE "$db"'
            WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$db')\gexec
EOSQL
    fi
done
IFS=$OLD_IFS

echo "Database initialization complete"
