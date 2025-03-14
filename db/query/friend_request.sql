-- name: GetFriendRequest :one
SELECT * FROM friend_requests 
WHERE 
	((sender_id = $1 AND receiver_id = $2) OR 
	(sender_id = $2 AND receiver_id = $1)) AND status = 1;

-- name: CreateFriendRequest :exec
INSERT INTO friend_requests (
    sender_id, receiver_id, request_desc
) VALUES (
    $1, $2, $3
);

-- name: UpdateFriendRequest :exec
UPDATE friend_requests
SET
	status  = $3,
	updated_at = now()
WHERE
sender_id = $1 AND receiver_id = $2 AND status = 1;