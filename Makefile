DB_URL=postgres://postgres:postgres@localhost:5432/dandelion?sslmode=disable

run:
	go run cmd/main.go

migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_URL}" -verbose down 1

migrateupall:
	migrate -path db/migration -database "${DB_URL}" -verbose up

migratedownall:
	migrate -path db/migration -database "${DB_URL}" -verbose down

newmigrate:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

.PHONLY: run migrateup migratedown migrateupall migratedownall newmigrate sqlc test