postgres:
	docker run --name postgres14 -p 3000:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres14 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:3000/simple_bank?sslmode=disable"  -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:3000/simple_bank?sslmode=disable"  -verbose down 
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/MasonPhan2110/SimpleBank/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown test server