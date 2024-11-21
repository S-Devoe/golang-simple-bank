createdb:
	docker exec -it postgres17 createdb -U root -O root simple_bank
# meaning of the above command: docker execute -iterative container-name createdb -username=root -owner=root database-name

dropdb:
	docker exec -it postgres17 dropdb simple_bank

postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

sqlc:
	sqlc generate

test: 
	@go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/S-Devoe/golang-simple-bank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server migratedown1 migrateup1 mock