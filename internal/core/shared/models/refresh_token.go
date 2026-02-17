package models

import "time"

type RefreshToken struct {
	ID        ID
	UserID    ID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRefreshToken(userID ID, tokenHash string, expiryDuration time.Duration) *RefreshToken {
	now := time.Now()
	return &RefreshToken{
		ID:        NewID(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(expiryDuration),
		RevokedAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked()
}
