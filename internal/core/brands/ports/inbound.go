package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	BrandService interface {
		List(ctx context.Context) ([]*models.Brand, error)
		DetailByID(ctx context.Context, id string) (*models.Brand, error)
		DetailBySlug(ctx context.Context, slug string) (*models.Brand, error)
		Update(ctx context.Context, id, name string) (*models.Brand, error)
		Create(ctx context.Context, name string) (*models.Brand, error)
		Delete(ctx context.Context, id string) error
	}
)
