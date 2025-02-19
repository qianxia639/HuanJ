-- name: ExistsFriendRequest :one
SELECT COUNT(*) FROM friend_requests 
WHERE 
	((user_id = $1 AND friend_id = $2) OR 
	(user_id = $2 AND friend_id = $1)) AND status = 1;

-- name: CreateFriendRequest :exec
INSERT INTO friend_requests (
    user_id, friend_id, request_desc
) VALUES (
    $1, $2, $3
);

-- name: UpdateFriendRequest :exec
UPDATE friend_requests
SET
	status  = $3,
	updated_at = now()
WHERE
user_id = $1 AND friend_id = $2 AND status = 1;