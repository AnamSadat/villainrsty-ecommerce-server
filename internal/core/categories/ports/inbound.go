package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	CategoryService interface {
		List(ctx context.Context) ([]*models.Category, error)
		DetailByID(ctx context.Context, id string) (*models.Category, error)
		DetailBySlug(ctx context.Context, slug string) (*models.Category, error)
		Update(ctx context.Context, id, name string) (*models.Category, error)
		Create(ctx context.Context, name string) (*models.Category, error)
		Delete(ctx context.Context, id string) error
	}
)
