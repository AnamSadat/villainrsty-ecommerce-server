package repository

import (
	"context"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/mapper"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type PasswordResetTokenRepository struct {
	q *sqlc.Queries
}

func NewPasswordResetTokenRepository(q *sqlc.Queries) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{q: q}
}

func (r *PasswordResetTokenRepository) Save(ctx context.Context, t *models.PasswordResetToken) error {
	var used pgtype.Timestamp
	if t.UsedAt != nil {
		used = pgtype.Timestamp{Time: *t.UsedAt, Valid: true}
	} else {
		used = pgtype.Timestamp{Valid: false}
	}

	return r.q.CreatePasswordResetToken(ctx, sqlc.CreatePasswordResetTokenParams{
		ID:        t.ID.String(),
		UserID:    t.UserID.String(),
		TokenHash: t.TokenHash,
		ExpiresAt: pgtype.Timestamp{Time: t.ExpiresAt, Valid: true},
		UsedAt:    used,
		CreatedAt: pgtype.Timestamp{Time: t.CreatedAt, Valid: true},
	})
}

func (r *PasswordResetTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	row, err := r.q.GetPasswordResetTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	return mapper.SQLPasswordResetTokenToDomain(row)
}

func (r *PasswordResetTokenRepository) MarkUsed(ctx context.Context, id models.ID) error {
	return r.q.MarkPasswordResetTokenUsed(ctx, id.String())
}
