## Build dependencies

 * Golang ^1.17 
 * swag ^1.8.2

## Runtime dependencies

 * Mongodb
 * Postgres | Mysql | Sqlite

## Bootstrap

Create .env file by .evn.example 

Run flags:

 * ```-v``` Verbose mode
 * ```-demo``` Fill database demo data
 * ```-docs``` Disable swagger public docs
 * ```-sentry``` Disable sentry
 * ```-mongo``` Disable mongo start initialization

At firs run use ```-m``` flag to create database structure

## API
### Auth
Set header ```db-key``` in each request. In socket methods set key in path.

### Api methods

See swagger in [docs dir](/docs)

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

## CLI

### Commands

#### admin
Create admin user
##### Usage:
```cli admin -e=admin.email.com -p=password```
##### Params:
* ```e``` Admin email
* ```p``` Admin password

#### backup
Backup databases
##### Usage:
```cli backup```
##### Optional params:
* ```s3``` Upload backup to s3
* 
#### restore
Restore databases from files
##### Usage:
```cli restore -dbPath=db.tar.gz -mongoPath=mongo.targ.gz```
##### Params:
* ```dbPath``` Path to bd backup
* ```mongoPath``` Path to mongo backup

#### migrate
Run db migrations
##### Usage:
```cli imgrate```

#### demo
Fill data base demo data
##### Usage:
```cli demo```
