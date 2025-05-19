#!/bin/sh

postgres_prefixes=("POSTGRES_USER" "DB_PASSWORD" "POSTGRES_DB" "POSTGRES_PORT" "REDIS_PORT" "KINOLK_AVATARS_FOLDER" "KINOLK_AVATARS_STATIC_PATH")

if [ -z "$1" ]; then
  echo "Usage: $0 <env_file>"
  echo "  <env_file>: Path to the file containing environment variables"
  exit 1
fi

env_file="$1"

if [ ! -f "$env_file" ]; then
  echo "Error: Environment file '$env_file' not found."
  exit 1
fi

is_postgres_var() {
  local var_name="$1"
  local prefix
  for prefix in "${postgres_prefixes[@]}"; do
    if [[ "${var_name}" == "${prefix}"* ]]; then
      return 0
    fi
  done
  return 1
}


while IFS= read -r line; do
# while IFS='=' read -r var_name var_value; do
  if [[ "$line" =~ ^# || -z "$line" ]]; then
    continue
  fi

  # Find the first equals sign
  IFS='=' read -r var_name var_value <<< "$line"

  var_value="${var_value//\"/}"

  if is_postgres_var "$var_name"; then
    export "$var_name"="$var_value"
    echo "Exported: $var_name"
  fi
done < "$env_file"
