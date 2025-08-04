#!/bin/sh
set -e

echo "Waiting for Postgres to be ready..."
while ! nc -z db 5432; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done

CMD=$1
if [ -z "${CMD}" ]; then
  echo "CMD is missing"
  exit 1
fi

# Perform migration
if [ "${CMD}" == "api" ]; then
  echo "Start database migration"
  /app/smart-pack migrate up
  echo "Finished database migration"
fi

# Start the SmartPack application
echo "Starting SmartPack application"
echo "Executing: /app/smart-pack ${CMD}"
/app/smart-pack "${CMD}"
