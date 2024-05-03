migrateup:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrateuptest:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank_test?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank?sslmode=disable" -verbose down 1


migratedowntest:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank_test?sslmode=disable" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/pakojabi/simplebank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migratedown migrateuptest migratedowntest sqlc test server mock migratedown1 migrateup1
