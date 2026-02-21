CREATE TABLE IF NOT EXISTS two_factor_otps (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id VARCHAR(36) NOT NULL UNIQUE,
    code_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_two_factor_otps_challenge_id ON two_factor_otps(challenge_id);
CREATE INDEX IF NOT EXISTS idx_two_factor_otps_expires_at ON two_factor_otps(expires_at);