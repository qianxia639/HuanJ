-- name: GetGroup :one
SELECT * FROM groups WHERE group_name = $1 LIMIT 1;

-- name: CreateGroup :one
INSERT INTO groups (
	group_name, creator_id, description
) VALUES (
	$1, $2, $3
) RETURNING *;