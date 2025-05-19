#!/bin/bash

set -e

start_compose_service() {
  set -e
  local service_path=$1

  source ./export_env_vars.sh ".env"

  echo "Запуск docker-compose для $service_path"

  # Проверяем, существует ли docker-compose.yml
  if [ ! -f "$service_path/docker-compose.yml" ]; then
    echo "docker-compose.yml не найден в $service_path. Пропуск..."
    return
  fi

  pushd "$service_path" > /dev/null
  docker-compose down -v
  docker-compose up --build -d
  popd > /dev/null
}

echo "Сборка Docker образа БД проекта"
start_compose_service "internal/db/postgresql_filmlk"


source ./export_env_vars.sh ".env"
echo $KINOLK_AVATARS_FOLDER
echo "Сборка Docker образа для корня проекта"
docker-compose down -v
docker-compose up --build -d

start_compose_service "user_service"
start_compose_service "auth_service"
start_compose_service "movie_service"

echo "Проверяем статус всех контейнеров..."
docker ps

echo "Все сервисы запущены!"
