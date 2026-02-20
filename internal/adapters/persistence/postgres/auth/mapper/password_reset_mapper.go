package mapper

import (
	"time"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

func SQLPasswordResetTokenToDomain(t sqlc.PasswordResetToken) (*models.PasswordResetToken, error) {
	var usedAt *time.Time
	if t.UsedAt.Valid {
		usedAt = &t.UsedAt.Time
	}

	return &models.PasswordResetToken{
		ID:        models.ID(t.ID),
		UserID:    models.ID(t.UserID),
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt.Time,
		UsedAt:    usedAt,
		CreatedAt: t.CreatedAt.Time,
	}, nil
}
