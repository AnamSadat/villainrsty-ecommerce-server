package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	CategoryRepository interface {
		Save(ctx context.Context, category *models.Category) error
		GetCategories(ctx context.Context) ([]*models.Category, error)
		GetDetailCategoryByID(ctx context.Context, id string) (*models.Category, error)
		GetDetailCategoryBySlug(ctx context.Context, slug string) (*models.Category, error)
		UpdateCategory(ctx context.Context, id, name, slug string) error
		DeleteCategory(ctx context.Context, id string) error
		ExistCategory(ctx context.Context, slug string) (bool, error)
	}
)
