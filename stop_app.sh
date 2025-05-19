#!/bin/bash

source ./export_env_vars.sh ".env"

remove_flag=""

if [[ "$1" == "--remove" ]]; then
  remove_flag="-v"
fi


srvs_paths=("internal/db/postgresql_filmlk" "user_service" "auth_service" "movie_service" "./")

for service_path in "${srvs_paths[@]}"; do
  echo "Запуск docker-compose для $service_path"

  # Проверяем, существует ли docker-compose.yml
  if [ ! -f "$service_path/docker-compose.yml" ]; then
    echo "docker-compose.yml не найден в $service_path. Пропуск..."
    return 1
  fi

  pushd "$service_path" > /dev/null

  if [[ "$1" == "--remove" ]]; then
    docker-compose down -v
  else
    docker-compose down
  fi  
  popd > /dev/null
done

if [[ "$1" == "--remove" ]]; then
  docker volume prune -f
fi

echo "Проверяем статус всех контейнеров..."
docker ps

echo "Все сервисы остановлены!"
