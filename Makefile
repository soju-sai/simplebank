include .env
export

dockerup:
	docker run --name postgres16 -e POSTGRES_USER=${SIMPLEBANK_USERNAME} -e POSTGRES_PASSWORD=${SIMPLEBANK_PASSWORD} -p 5432:5432 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --owner=${SIMPLEBANK_USERNAME} --username=${SIMPLEBANK_USERNAME} simple_bank "A database acts as a bank db"

dropdb:
	docker exec -it postgres16 dropdb --username=${SIMPLEBANK_USERNAME} simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://${SIMPLEBANK_USERNAME}:${SIMPLEBANK_PASSWORD}@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://${SIMPLEBANK_USERNAME}:${SIMPLEBANK_PASSWORD}@localhost:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: dockerup createdb dropdb migrateup migratedown