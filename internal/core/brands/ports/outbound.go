package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	BrandRepository interface {
		Save(ctx context.Context, brand *models.Brand) error
		GetBrands(ctx context.Context) ([]*models.Brand, error)
		GetDetailBrandByID(ctx context.Context, id string) (*models.Brand, error)
		GetDetailBrandBySlug(ctx context.Context, slug string) (*models.Brand, error)
		UpdateBrand(ctx context.Context, id, name, slug string) error
		DeleteBrand(ctx context.Context, id string) error
		ExistBrand(ctx context.Context, slug string) (bool, error)
	}
)
