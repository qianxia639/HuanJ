-- name: CreateMessage :exec
INSERT INTO messages (
    session_id, sender_id, receiver_id, send_type, receiver_type, message_type, content
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);