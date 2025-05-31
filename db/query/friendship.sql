-- name: ExistsFriendship :one
SELECT EXISTS(
    SELECT 1 FROM friend_requests
    WHERE (from_user_id = $1 AND to_user_id = $2 AND status = 2)
    OR (from_user_id = $2 AND to_user_id = $1 AND status = 2)
);

-- name: CreateFriendship :one
INSERT INTO friendships (
    user_id, friend_id, note
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetFriendList :many
SELECT * FROM friendships WHERE user_id = $1;

-- name: DeleteFriend :exec
DELETE FROM friendships 
WHERE (user_id = $1 AND friend_id = $2) 
    OR (user_id = $2 AND friend_id = $1);