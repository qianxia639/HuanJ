-- name: CreateGroupMember :one
INSERT INTO group_members (
	group_id, user_id, role, agreed
) VALUES (
	$1, $2, $3, $4
)
RETURNING *;