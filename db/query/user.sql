-- name: CreateUser :one
INSERT INTO users (
	username, nickname, password, email, gender
) VALUES (
	$1, $2, $3, $4, $5
)
RETURNING *;