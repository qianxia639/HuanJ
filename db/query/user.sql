-- name: CreateUser :one
INSERT INTO users (
	username, nickname, password, email, gender
) VALUES (
	$1, $2, $3, $4, $5
)
RETURNING *;

-- name: ExistsUsername :one
SELECT EXISTS (
	SELECT 1 FROM users 
	WHERE username = $1
);

-- name: ExistsEmail :one
SELECT EXISTS (
	SELECT 1 FROM users
	WHERE email = $1
)

-- name: ExistsNickname :one
SELECT EXISTS (
	SELECT 1 FROM users
	WHERE nickname = $1
);

-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users 
SET
	gender = COALESCE(sqlc.narg('gender'), gender), 
	nickname = COALESCE(sqlc.narg('nickname'), nickname), 
	updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdatePwd :exec
UPDATE users
SET
	password = $2,
	password_changed_at = now(),
	updated_at = now()
WHERE
	email = $1;
