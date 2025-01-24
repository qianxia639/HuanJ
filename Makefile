DB_URL=postgres://postgres:postgres@localhost:5432/dandelion?sslmode=disable

server:
	go run main.go

migrateup:
	migrate -path internal/db/migration -database "${DB_URL}" -verbose up

migratedown:
	migrate -path internal/db/migration -database "${DB_URL}" -verbose down

migrateupall:
	migrate -path internal/db/migration -database "${DB_URL}" -verbose up 1

migratedownall:
	migrate -path internal/db/migration -database "${DB_URL}" -verbose down 1

newmigrate:
	migrate create -ext sql -dir db/migration -seq $(name)

git:
	git filter-branch --force --index-filter 'git rm --cached --ignore-unmatch internal/config/config.toml' --prune-empty --tag-name-filter cat -- --all

.PHONLY: server migrateup migratedown migrateupall migratedownall newmigrate git