package repository

import (
	"context"
	"errors"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/mapper"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	appErr "villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TwoFactorOTPRepository struct {
	q *sqlc.Queries
}

func NewTwoFactorOTPRepository(q *sqlc.Queries) *TwoFactorOTPRepository {
	return &TwoFactorOTPRepository{q: q}
}

func (r *TwoFactorOTPRepository) Save(ctx context.Context, otp *models.TwoFactorOTP) error {
	var usedAt pgtype.Timestamp
	if otp.UsedAt != nil {
		usedAt = pgtype.Timestamp{Time: *otp.UsedAt, Valid: true}
	}

	return r.q.CreateTwoFactorOTP(ctx, sqlc.CreateTwoFactorOTPParams{
		ID:          otp.ID.String(),
		UserID:      otp.UserID.String(),
		ChallengeID: otp.ChallengeID,
		CodeHash:    otp.CodeHash,
		ExpiresAt:   pgtype.Timestamp{Time: otp.ExpiresAt, Valid: true},
		UsedAt:      usedAt,
		CreatedAt:   pgtype.Timestamp{Time: otp.CreatedAt, Valid: true},
	})
}

func (r *TwoFactorOTPRepository) GetByChallengeID(ctx context.Context, challengeID string) (*models.TwoFactorOTP, error) {
	row, err := r.q.GetTwoFactorOTPByChallengeID(ctx, challengeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "2fa challenge not found")
		}
		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get 2fa challenge", err)
	}

	return mapper.SQLTwoFactorOTPToDomain(row), nil
}

func (r *TwoFactorOTPRepository) MarkUsed(ctx context.Context, id models.ID) error {
	return r.q.MarkTwoFactorOTPUsed(ctx, id.String())
}

func (r *TwoFactorOTPRepository) DeleteExpired(ctx context.Context) error {
	return r.q.DeleteExpiresTwoFactorOTP(ctx)
}
