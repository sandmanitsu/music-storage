## Тестовое задание Реализация онлайн библиотеки песен
```Структура проекта```
- cmd/app - точка входа в приложение
- docs - сгенерированные swagger файлы
- migrations - файлы для миграции

```API```
- /api/v1/track/list - получение списка песни
- /api/v1/track/text - получение текста песни
- /api/v1/track/delete - удаление песни
- /api/v1/track/update - обновление данных песни

##
- /swagger/ - swaggerUI
##

## .env file
```bash
HOST=хост
PORT=порт для приложения
ENV=среда работы приложения # dev, prod

DBHOST=хост БД
USER=логин БД
PASSWORD=пароль БД
DBNAME=название БД
DBPORT=порт БД

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://$USER:$PASSWORD@$DBHOST:$DBPORT/$DBNAME
GOOSE_MIGRATION_DIR=./migrations
```

##  Makefile

```bash
# запуск приложения
run: 
	go run cmd/app/main.go
# миграции:
up:
	goose up

down:
	goose down
# генерация swagger:
swag:
	swag init -g cmd/app/main.go
```
