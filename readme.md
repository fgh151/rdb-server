## Зависимости

Mongodb

Postgres | Mysql | Sqlite

## Запуск

Создать .env файл см .env.example

Флаги запуска:

 * ```-v``` режим подробного логирования
 * ```-m``` запуск миграций при старте
 * ```-demo``` заполнение баз данных демо данными

## API
### Авторизация
В каждом запросе надо указывать заголовок 
db-key в котором указывать ключ

В методе подписки ключ нужно казать не в заголовке, а в параметре адреса

### Доступные методы

POST /push/{topic} Отправка сообщения в топик

PATCH /push/{topic} Изменение сообщения

POST /find/{topic} Поиск сообщений

GET /list/{topic} Список сообщений

GET /subscribe/{topic}/{key} Сокет для подписки на новые сообщения

## Docker compose example

```yaml
version: '3.7'

services:
  mongodb:
    restart: always
    image: mongo:4.4.0-bionic
    container_name: mongodb
    environment:
      MONGO_INITDB_ROOT_USERNAME: dockerMongoAdmin
      MONGO_INITDB_ROOT_PASSWORD: dockerMongoPassword
      MONGO_INITDB_DATABASE: dockerdb
    ports:
      - "27017:27017"
    volumes:
      - .docker/mongo-init.js:/docker-entrypoint-initdb.d/mongo-init.js:ro

  app:
    restart: always
    image: fgh151/db-server:0.0.1
    depends_on:
      - mongodb
    container_name: db-server
    ports:
      - "9090:9090"
    env_file:
      - .env
    volumes:
      - .env:/.env
      - ./db.db:/db.db
    links:
      - mongodb
  
  watchtower:
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```