# PostgreSQL Migrations with golang-migrate

This guide explains how to create and run database migrations for the relayer services using `golang-migrate`.

## Installation

Install the `golang-migrate` tool:

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

**Mac:**
```bash
brew install golang-migrate
```

## Creating Migration Files

Create new migration files with SQL extension:

```bash
migrate create -ext sql -dir private-relayer/repository/migrations -seq migration_name
```

This creates two files:
- `XXXXXX_migration_name.up.sql` - Applied when migrating up
- `XXXXXX_migration_name.down.sql` - Applied when migrating down

## Running Migrations

### Using Docker

**Migration Up:**
```bash
docker run -v $(pwd)/private-relayer/repository/migrations:/migrations --network host migrate/migrate \
  -path=/migrations/ \
  -database "postgres://USER:PASSWORD@localhost:5432/DATABASE?sslmode=disable" \
  up
```

**Migration Down:**
```bash
docker run -v $(pwd)/private-relayer/repository/migrations:/migrations --network host migrate/migrate \
  -path=/migrations/ \
  -database "postgres://USER:PASSWORD@localhost:5432/DATABASE?sslmode=disable" \
  down
```

### Using Local CLI

**Migration Up:**
```bash
migrate -path private-relayer/repository/migrations \
  -database "postgres://USER:PASSWORD@localhost:5432/DATABASE?sslmode=disable" \
  up
```

**Migration Down:**
```bash
migrate -path private-relayer/repository/migrations \
  -database "postgres://USER:PASSWORD@localhost:5432/DATABASE?sslmode=disable" \
  down
```

### Local Development Example

For the local Docker development environment:

```bash
migrate -path private-relayer/repository/migrations \
  -database "postgres://admin:admin@localhost:5432/relayer_a?sslmode=disable" \
  up
```

## Embedded Migrations

Note: In production, migrations are embedded in the Go binary using `//go:embed` and run automatically on startup. Manual migration is typically only needed for:
- Local development troubleshooting
- Creating new migrations
- Manual database maintenance

## Useful Commands

**Check current migration version:**
```bash
migrate -path private-relayer/repository/migrations \
  -database "postgres://admin:admin@localhost:5432/relayer_a?sslmode=disable" \
  version
```

**Force a specific version (use with caution):**
```bash
migrate -path private-relayer/repository/migrations \
  -database "postgres://admin:admin@localhost:5432/relayer_a?sslmode=disable" \
  force VERSION
```

**Migrate to a specific version:**
```bash
migrate -path private-relayer/repository/migrations \
  -database "postgres://admin:admin@localhost:5432/relayer_a?sslmode=disable" \
  goto VERSION
```
