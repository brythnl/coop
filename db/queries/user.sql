-- name: CreateUser :one
INSERT INTO users (
	username,
	password_hash
) VALUES (
	$1, $2
)
RETURNING id, username, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT id, username, password_hash
FROM users
WHERE username = $1
LIMIT 1;


