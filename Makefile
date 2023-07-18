include .env

## Delete app
app-rm:
	docker rm -f app-container
	docker rmi app-image

## Build and run app
app-up:
	docker build -t app-image .
	docker run -dp 8080:8080\
		--network order-and-pay-app --network-alias app\
		-e POSTGRES_HOST=${POSTGRES_HOST}\
		-e POSTGRES_USER=${POSTGRES_USER}\
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD}\
		-e POSTGRES_DB=${POSTGRES_DB}\
		-e POSTGRES_PORT=${POSTGRES_PORT}\
		--name app-container\
		app-image

#========================#
#== DATABASE MIGRATION ==#
#========================#

## Run migrations UP
migrate-up:
	docker compose run --rm migrate up

## Run migrations DOWN
migrate-down:
	docker compose run --rm migrate down 1

## Create a DB migration files e.g `make migrate-create name=migration-name`
migrate-create:
	docker compose run --rm migrate create -ext sql -dir /migrations -seq $(name)

## Enter to database console
shell-db:
	docker compose exec postgresql psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}