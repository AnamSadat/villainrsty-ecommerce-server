package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	ProductRepository interface {
		Save(ctx context.Context, product *models.Product) error
		GetProducts(ctx context.Context) ([]*models.Product, error)
		GetDetailProductByID(ctx context.Context, id string) (*models.Product, error)
		GetDetailProductBySlug(ctx context.Context, slug string) (*models.Product, error)
		UpdateProduct(ctx context.Context, id, brandID, categoryID, name, slug, description string, isActive bool) error
		DeleteProduct(ctx context.Context, id string) error
		ExistProduct(ctx context.Context, name string) (bool, error)
	}
)
