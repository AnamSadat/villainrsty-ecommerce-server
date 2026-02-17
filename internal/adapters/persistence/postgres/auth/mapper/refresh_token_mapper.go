package mapper

import (
	"time"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func SQLRefreshTokenToDomain(t sqlc.RefreshToken) *models.RefreshToken {
	var revokedAt *time.Time
	if t.RevokedAt.Valid {
		revokedAt = &t.RevokedAt.Time
	}

	return &models.RefreshToken{
		ID:        models.ID(t.ID),
		UserID:    models.ID(t.UserID),
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt.Time,
		RevokedAt: revokedAt,
		CreatedAt: t.CreatedAt.Time,
		UpdatedAt: t.UpdatedAt.Time,
	}
}

func DomainRefreshTokenToSQLCParams(t *models.RefreshToken) sqlc.CreateRefreshTokenParams {
	return sqlc.CreateRefreshTokenParams{
		ID:        t.ID.String(),
		UserID:    t.UserID.String(),
		TokenHash: t.TokenHash,
		ExpiresAt: pgtype.Timestamp{
			Time:  t.ExpiresAt,
			Valid: true,
		},
		RevokedAt: pgtype.Timestamp{
			Time:  time.Time{},
			Valid: false,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  t.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  t.UpdatedAt,
			Valid: true,
		},
	}
}
