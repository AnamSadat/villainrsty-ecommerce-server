-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, used_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetPasswordResetTokenByHash :one 
SELECT id, user_id, token_hash, expires_at, used_at, created_at
FROM password_reset_tokens
WHERE token_hash = $1
LIMIT 1;

-- name: MarkPasswordResetTokenUsed :exec
UPDATE password_reset_tokens
SET used_at = NOW()
WHERE id = $1;

-- name: DeleteExpirePasswordResetToken :exec
DELETE FROM password_reset_tokens
WHERE expires_at < NOW() OR used_at IS NOT NULL;