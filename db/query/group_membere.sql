-- name: CreateGroupMember :one
INSERT INTO group_members (
	group_id, user_id, role
) VALUES (
	$1, $2, $3
)
RETURNING *;

-- name: ExistsGroupMember :one
SELECT EXISTS (
	SELECT 1 FROM group_members
	WHERE group_id = $1 AND user_id = $2
);

-- name: GetGroupMemberList :many
SELECT * FROM group_members
WHERE group_id = $1;