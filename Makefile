dockerup:
	docker run --name postgres16 -e POSTGRES_USER=${DB_USERNAME} -e POSTGRES_PASSWORD=${DB_PASSWORD} -p 5432:5432 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --owner=${DB_USERNAME} --username=${DB_USERNAME} simple_bank "A database acts as a bank db"

dropdb:
	docker exec -it postgres16 dropdb --username=${DB_USERNAME} simple_bank

migrateup:
	migrate -path db/migration -database ${DB_SOURCE} -verbose up

migratedown:
	migrate -path db/migration -database ${DB_SOURCE} -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mockdb:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store

.PHONY: dockerup createdb dropdb migrateup migratedown sqlc test server mockdb