package mapper

import (
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func SQLBrandBySlugToDomain(b sqlc.Brand) *models.Brand {
	return &models.Brand{
		ID:        models.ID(b.ID),
		Name:      b.Name,
		Slug:      b.Slug,
		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
	}
}

func DomainBrandToSQLParams(b *models.Brand) sqlc.CreateBrandParams {
	return sqlc.CreateBrandParams{
		ID:   b.ID.String(),
		Name: b.Name,
		Slug: b.Slug,
		CreatedAt: pgtype.Timestamp{
			Time:  b.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  b.UpdatedAt,
			Valid: true,
		},
	}
}
