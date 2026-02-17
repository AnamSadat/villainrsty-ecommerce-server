package models

import "time"

type PasswordResetToken struct {
	ID        ID
	UserID    ID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (t *PasswordResetToken) IsValid() bool {
	if t.UsedAt != nil {
		return false
	}
	return time.Now().Before(t.ExpiresAt)
}
