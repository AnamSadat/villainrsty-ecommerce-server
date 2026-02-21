package models

import "time"

type TwoFactorOTP struct {
	ID          ID
	UserID      ID
	ChallengeID string
	CodeHash    string
	ExpiresAt   time.Time
	UsedAt      *time.Time
	CreatedAt   time.Time
}

func (o *TwoFactorOTP) IsValid() bool {
	if o.UsedAt != nil {
		return false
	}

	return time.Now().Before(o.ExpiresAt)
}
