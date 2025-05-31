-- name: CreateUser :one
INSERT INTO users (
	username, nickname, password, email, gender
) VALUES (
	$1, $2, $3, $4, $5
)
RETURNING *;

-- name: ExistsUsername :one
SELECT COUNT(*) FROM users WHERE username = $1;

-- name: ExistsEmail :one
SELECT COUNT(*) FROM users WHERE email = $1;

-- name: ExistsNickname :one
SELECT COUNT(*) FROM users WHERE nickname = $1;;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users 
SET 
	gender = $1, 
	nickname = $2, 
	updated_at = now()
WHERE id = $3
AND (gender IS DISTINCT FROM $1 OR nickname IS DISTINCT FROM $2);

-- name: UpdatePwd :exec
UPDATE users
SET
	password = $1,
	password_changed_at = now(),
	updated_at = now()
WHERE
	id = $2 AND email = $3;