package mapper

import (
	"time"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

func SQLTwoFactorOTPToDomain(t sqlc.TwoFactorOtp) *models.TwoFactorOTP {
	var usedAt *time.Time
	if t.UsedAt.Valid {
		usedAt = &t.UsedAt.Time
	}

	return &models.TwoFactorOTP{
		ID:          models.ID(t.ID),
		UserID:      models.ID(t.UserID),
		ChallengeID: t.ChallengeID,
		CodeHash:    t.CodeHash,
		ExpiresAt:   t.ExpiresAt.Time,
		UsedAt:      usedAt,
		CreatedAt:   t.CreatedAt.Time,
	}
}
