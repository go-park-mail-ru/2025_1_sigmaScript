# USAGE

## Переходим в директорию с постгрес
```bash
cd db/postgresql_filmlk/
```

## Экспортируем перепенные среды постгрес из <env_file>
```bash
source ../export_env_vars.sh <env_file>
```
## Поднимаем PostgreSQL с помощью docker compose
```bash
docker-compose up
```
или как демон
```bash
 docker-compose up -d
```
## Убиваем PostgreSQL с помощью docker compose (если запустили с -d)
```bash
docker-compose stop
```

## Удаляем контейнер PostgreSQL с помощью docker compose (если запустили с -d)
```bash
docker-compose down
```

## Подключение к контейнеру в постгрес (если надо что-то руками изменить)
```bash
docker exec -it <docker_postgres_container> psql -U <POSTGRES_USER> -d POSTGRES_DB
```

