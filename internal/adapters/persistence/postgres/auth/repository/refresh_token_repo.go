package repository

import (
	"context"
	"errors"
	"time"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/mapper"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	appErr "villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
	"villainrsty-ecommerce-server/pkg/validator"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshTokenRepository struct {
	queris    *sqlc.Queries
	validator *validator.Validator
}

func NewRefreshTokenRepository(queries *sqlc.Queries) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		queris:    queries,
		validator: validator.NewValidate(),
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, t *models.RefreshToken) error {
	if err := r.validator.ValidateRequired("token_hash", t.TokenHash); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid_token_hash", err)
	}

	if err := r.validator.ValidateRequired("user_id", t.UserID.String()); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid user id", err)
	}

	params := mapper.DomainRefreshTokenToSQLCParams(t)
	if err := r.queris.CreateRefreshToken(ctx, params); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to create refresh token", err)
	}

	return nil
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	if err := r.validator.ValidateRequired("token_hash", tokenHash); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invalid token hash", err)
	}

	row, err := r.queris.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "refresh token not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get refresh token", err)
	}

	token := mapper.SQLRefreshTokenToDomain(row)

	return token, nil
}

func (r *RefreshTokenRepository) GetByUserID(ctx context.Context, userID models.ID) ([]*models.RefreshToken, error) {
	if err := r.validator.ValidateRequired("user_id", userID.String()); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invalid user id", err)
	}

	rows, err := r.queris.GetRefreshTokensByUserID(ctx, userID.String())
	if err != nil {
		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get refresh token", err)
	}

	tokens := make([]*models.RefreshToken, len(rows))
	for i, row := range rows {
		tokens[i] = mapper.SQLRefreshTokenToDomain(row)
	}

	return tokens, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenID models.ID) error {
	if err := r.validator.ValidateRequired("token_id", tokenID.String()); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid token id", err)
	}

	now := time.Now()
	if err := r.queris.RevokeRefreshToken(ctx, sqlc.RevokeRefreshTokenParams{
		ID:        tokenID.String(),
		RevokedAt: pgtype.Timestamp{Time: now, Valid: true},
	}); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to revoke refresh token", err)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	if err := r.queris.DeleteExpiredRefreshTokens(ctx, pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to delete expired refresh tokens", err)
	}

	return nil
}
