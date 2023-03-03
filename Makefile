postgres:
	docker run --name postgres14 --network bank-network -p 3000:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres14 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:e1fdFRekMRW1aYNrHZ0F@simple-bank.cenlrkxivm0u.ap-southeast-1.rds.amazonaws.com:5432/simple_bank"  -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:3000/simple_bank?sslmode=disable"  -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:3000/simple_bank?sslmode=disable"  -verbose down 
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:3000/simple_bank?sslmode=disable"  -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/MasonPhan2110/SimpleBank/db/sqlc Store
env:
	aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > .env
.PHONY:env postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 test server