#!/bin/sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
while ! pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "PostgreSQL is not ready yet. Waiting..."
  sleep 2
done

echo "PostgreSQL is ready. Running migrations..."

# Set PGPASSWORD for authentication
export PGPASSWORD="$DB_PASSWORD"

# Run migrations
for migration in /app/migrations/*.sql; do
  if [ -f "$migration" ]; then
    echo "Running migration: $(basename $migration)"
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration"
    if [ $? -eq 0 ]; then
      echo "Migration $(basename $migration) completed successfully"
    else
      echo "Migration $(basename $migration) failed"
      exit 1
    fi
  fi
done

echo "All migrations completed successfully!"
