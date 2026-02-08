package mapper

import (
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	// "villainrsty-ecommerce-server/internal/core/auth/models"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func SQLCUserByEmailToDomain(sqlcUser sqlc.GetUserByEmailRow) *models.User {
	return &models.User{
		ID:        models.ID(sqlcUser.ID),
		Email:     sqlcUser.Email,
		Password:  sqlcUser.Password,
		Name:      sqlcUser.Name,
		CreatedAt: sqlcUser.CreatedAt.Time,
		UpdatedAt: sqlcUser.UpdatedAt.Time,
	}
}

func SQLCUserByIDToDomain(sqlcUser sqlc.GetUserByIDRow) *models.User {
	return &models.User{
		ID:        models.ID(sqlcUser.ID),
		Email:     sqlcUser.Email,
		Password:  sqlcUser.Password,
		Name:      sqlcUser.Name,
		CreatedAt: sqlcUser.CreatedAt.Time,
		UpdatedAt: sqlcUser.UpdatedAt.Time,
	}
}

func DomainUserToSQLCParams(user *models.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		CreatedAt: pgtype.Timestamp{
			Time:  user.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  user.UpdatedAt,
			Valid: true,
		},
	}
}
