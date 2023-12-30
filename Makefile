CONTAINER := $(or $(POSTGRES_CONTAINER), some-postgres)

postgres:
	docker run --name $(CONTAINER) -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres123 -d postgres:12-alpine


createdb:
	docker exec -it $(CONTAINER) createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it $(CONTAINER) dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres123@localhost:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: createdb dropdb postgres migrateup
