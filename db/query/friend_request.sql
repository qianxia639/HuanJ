-- 检查申请是否存在
-- name: GetFriendRequest :one
-- SELECT * FROM friend_requests 
-- WHERE 
-- 	from_user_id = $1 AND to_user_id = $2 AND status = 1;
SELECT EXISTS (
	SELECT 1
	FROM friend_requests
	WHERE from_user_id = $1 AND to_user_id = $2 AND status = 1
);

-- name: CreateFriendRequest :exec
INSERT INTO friend_requests (
    from_user_id, to_user_id, request_desc
) VALUES (
    $1, $2, $3
);

-- name: UpdateFriendRequest :exec
UPDATE friend_requests
SET
	status  = $3,
	updated_at = now()
WHERE
	(from_user_id = $1 AND to_user_id = $2 AND status = 1)
OR 
	(to_user_id = $1 AND from_user_id = $2 AND status = 1);

-- name: ListFriendRequestByPending :many
SELECT * FROM friend_requests 
WHERE 
	to_user_id = $1 AND status = 1
ORDER BY created_at DESC;