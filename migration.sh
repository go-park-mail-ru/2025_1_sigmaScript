#!/bin/sh
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
NAME=$1
DIR="internal/db/postgresql_filmlk/migrations"
EXTENSION="sql"

if [ -z "$NAME" ]; then
  echo "Usage: $0 <migration_name>"
  exit 1
fi

mkdir -p "$DIR"

UP_FILE="${DIR}/${TIMESTAMP}_${NAME}.up.${EXTENSION}"
DOWN_FILE="${DIR}/${TIMESTAMP}_${NAME}.down.${EXTENSION}"

touch "$UP_FILE" "$DOWN_FILE"

echo "Migration created:"
echo "  $UP_FILE"
echo "  $DOWN_FILE"
