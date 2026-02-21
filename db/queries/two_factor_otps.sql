-- name: CreateTwoFactorOTP :exec
INSERT INTO two_factor_otps (id, user_id, challenge_id, code_hash, expires_at, used_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetTwoFactorOTPByChallengeID :one
SELECT id, user_id, challenge_id, code_hash, expires_at, used_at, created_at
FROM two_factor_otps
WHERE challenge_id = $1
LIMIT 1;

-- name: MarkTwoFactorOTPUsed :exec
UPDATE two_factor_otps
SET used_at = NOW()
WHERE id = $1;

-- name: DeleteExpiresTwoFactorOTP :exec
DELETE FROM two_factor_otps
WHERE expires_at < NOW() OR used_at IS NOT NULL;