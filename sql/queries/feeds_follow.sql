-- name: CreateFeedFollow :one
INSERT INTO feed_follow (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
--

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follow WHERE feed_id = $1 AND user_id = $2;
--
