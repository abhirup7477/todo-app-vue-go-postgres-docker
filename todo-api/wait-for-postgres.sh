#!/bin/sh
# wait-for-postgres.sh
set -e

host="$1"
shift
cmd="$@"

echo "Waiting for Postgres at $host..."
until pg_isready -h "$host" -U "admin"; do
  sleep 2
done

echo "Postgres is up - executing command"
exec $cmd