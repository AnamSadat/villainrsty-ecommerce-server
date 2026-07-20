package mapper

import (
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func SQLCategoryBySlugToDomain(c sqlc.Category) *models.Category {
	return &models.Category{
		ID:        models.ID(c.ID),
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}

func DomainCategoryToSQLParams(c *models.Category) sqlc.CreateCategoryParams {
	return sqlc.CreateCategoryParams{
		ID:   c.ID.String(),
		Name: c.Name,
		Slug: c.Slug,
		CreatedAt: pgtype.Timestamp{
			Time:  c.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  c.UpdatedAt,
			Valid: true,
		},
	}
}
