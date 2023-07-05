include .env

##Build and Run app container
up-app:
	docker build -t app .
	docker run -dp 8080:8080\ 
	-e POSTGRES_HOST=${POSTGRES_HOST}\
	-e POSTGRES_USER=${POSTGRES_USER}\
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD}\
	-e POSTGRES_DB=${POSTGRES_DB}\
	-e POSTGRES_PORT=${POSTGRES_PORT}\
	app

#========================#
#== DATABASE MIGRATION ==#
#========================#

## Run migrations UP
migrate-up:
	docker compose run --rm migrate up

## Rollback migrations against non test DB
migrate-down:
	docker compose run --rm migrate down 1

## Create a DB migration files e.g `make migrate-create name=migration-name`
migrate-create:
	docker compose --rm migrate create -ext sql -dir /migrations -seq $(name)

## Enter to database console
shell-db:
	docker compose exec postgresql psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}