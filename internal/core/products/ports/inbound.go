package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	ProductService interface {
		List(ctx context.Context) ([]*models.Product, error)
		DetailByID(ctx context.Context, id string) (*models.Product, error)
		DetailBySlug(ctx context.Context, slug string) (*models.Product, error)
		Create(ctx context.Context, id, brandID, categoryID, name, slug, description string, isActive bool) (*models.Product, error)
		Update(ctx context.Context, id, brandID, categoryID, name, slug, description string, isActive bool) (*models.Product, error)
		Delete(ctx context.Context, id string) error
	}
)
