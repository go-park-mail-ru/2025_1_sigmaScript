#!/bin/bash


build_flag=""

if [[ "$1" == "--build" ]]; then
  build_flag="$1"
fi

source ./export_env_vars.sh ".env"

srvs_paths=("internal/db/postgresql_filmlk" "user_service" "auth_service" "movie_service" "./")

for service_path in "${srvs_paths[@]}"; do
  echo "Запуск docker-compose для $service_path"

  # Проверяем, существует ли docker-compose.yml
  if [ ! -f "$service_path/docker-compose.yml" ]; then
    echo "docker-compose.yml не найден в $service_path. Пропуск..."
    return 1
  fi

  pushd "$service_path" > /dev/null
  if [[ "$1" == "--build" ]]; then
    echo "building docker-compose for $service_path..."
    docker-compose down -v
    docker-compose up --build -d
  else
    docker-compose down -v
    docker-compose up -d
  fi
  popd > /dev/null
done

echo "Проверяем статус всех контейнеров..."
docker ps

echo "Все сервисы запущены!"
