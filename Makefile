DB_URL=postgres://postgres:postgres@localhost:5432/dandelion?sslmode=disable

server:
	go run main.go

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

.PHONLY: server migrateup migratedown migrateupall migratedownall newmigrate