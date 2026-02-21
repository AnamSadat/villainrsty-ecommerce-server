package repository

import (
	"context"
	"errors"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	appErr "villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5"
)

type UserRoleRepository struct {
	q *sqlc.Queries
}

func NewUserRoleRepository(q *sqlc.Queries) *UserRoleRepository {
	return &UserRoleRepository{q: q}
}

func (r *UserRoleRepository) AssignRole(ctx context.Context, userID models.ID, roleName string) error {
	roleID, err := r.q.GetRoleByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appErr.New(appErr.ErrNotFound, "role not found")
		}
		return appErr.Wrap(appErr.ErrInternal, "failed to get role id", err)
	}

	if err := r.q.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
		UserID: userID.String(),
		RoleID: roleID,
	}); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to assign role", err)
	}

	return nil
}

func (r *UserRoleRepository) GetPrimaryRoleByUser(ctx context.Context, userID models.ID) (string, error) {
	name, err := r.q.GetPrimaryRoleByUserID(ctx, userID.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", appErr.New(appErr.ErrNotFound, "role not found")
		}

		return "", appErr.Wrap(appErr.ErrInternal, "failed to get user role", err)
	}

	return name, nil
}
