package mapper

import (
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func SQLCUserToDomain(u sqlc.User) *models.User {
	return &models.User{
		ID:        models.ID(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		Role:      "customer",
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
}

func DomainUserToSQLCParams(u *models.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		ID:       u.ID.String(),
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		CreatedAt: pgtype.Timestamp{
			Time:  u.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  u.UpdatedAt,
			Valid: true,
		},
	}
}
