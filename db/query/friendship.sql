-- 是否已是好友关系
-- name: ExistsFriendship :one
SELECT EXISTS(
    SELECT 1 FROM friendships
    WHERE (user_id = $1 AND friend_id = $2)
);

-- name: CreateFriendship :one
INSERT INTO friendships (
    user_id, friend_id, note
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetFriendList :many
SELECT u.*
FROM friendships f
JOIN users u ON f.friend_id = u.id
WHERE f.user_id = $1;

-- name: DeleteFriend :exec
-- DELETE FROM friendships 
-- WHERE (user_id = $1 AND friend_id = $2) 
--     OR (user_id = $2 AND friend_id = $1);
DELETE FROM friendships 
WHERE (user_id, friend_id) 
IN (($1, $2), ($2, $1));