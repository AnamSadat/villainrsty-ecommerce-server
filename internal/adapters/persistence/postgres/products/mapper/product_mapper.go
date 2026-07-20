package mapper

import (
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func SQLProductToDamain(p sqlc.Product) *models.Product {
	return &models.Product{
		ID:          models.ID(p.ID),
		BrandID:     models.ID(p.BrandID),
		CategoryID:  models.ID(p.CategoryID),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}

func DomainProductToSQLParams(p *models.Product) sqlc.CreateProductParams {
	return sqlc.CreateProductParams{
		ID:          p.ID.String(),
		BrandID:     p.BrandID.String(),
		CategoryID:  p.CategoryID.String(),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		IsActive:    p.IsActive,
		CreatedAt: pgtype.Timestamp{
			Time:  p.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  p.UpdatedAt,
			Valid: true,
		},
	}
}
