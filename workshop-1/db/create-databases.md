# initdb.d/000-create-databases.sh

## Назначение / Purpose

Скрипт для автоматического создания баз данных PostgreSQL и пользователей на основе конфигурационного файла.  
Automates PostgreSQL database and user creation based on a configuration file.

## Требования (Requirements)
- PostgreSQL 9.2+ (используется синтаксис `\gexec`)
- Доступ с правами `postgres` пользователя  
  Requires `postgres` user privileges

## Конфигурация / Configuration

```sh
# config/databases.conf
db_list_to_create='cart order'  # Database IDs

db_cart_name='cart'
db_cart_user='cart_user'
db_cart_password='Str0ngPass!'

db_order_name='orders'
db_order_user='order_user'
db_order_password='Order@123'
```

## Интеграция с Docker / Docker Integration

```yaml
# docker-compose.yml
services:
  db:
    image: postgres
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - ./db/config:/config
      - ./db/initdb.d:/docker-entrypoint-initdb.d
      - ./db/.docker/postgresql/data:/var/lib/postgresql/data
```

## Использование / Usage

### Первичная настройка / Initial setup

При первом запуске контейнера (если в каталоге PGDATА отсутствуют данные) скрипт выполнится автоматически.  
When the container is started for the first time (if there is no data in the PGDATA directory), the script will run automatically.

### Добавление БД (без потери данных) / Add DB (no data loss)

1. Обновите конфиг:  
Update config:

```sh
db_list_to_create='cart order payment'  # Add new DB payment

db_payment_name='payment'
db_payment_user='payment_user'
db_payment_password='Pay@2023'
```

2. Выполните скрипт вручную:  
Execute the script manually:

```bash
docker-compose exec db sh -c '/docker-entrypoint-initdb.d/000-create-databases.sh'
```

## Ограничения / Limitations

- Скрипт выполняется автоматически только при первом запуске контейнера.  
The script is executed automatically only when the container is started for the first time.

- Скрипт не изменяет существующие базы данных и пользователей.  
The script does not modify existing databases and users.
