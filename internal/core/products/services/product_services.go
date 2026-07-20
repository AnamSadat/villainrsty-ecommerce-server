package services

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/products/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type ProductServices struct {
	productRepo ports.ProductRepository
}

func NewProductService(productRepo ports.ProductRepository) *ProductServices {
	return &ProductServices{
		productRepo: productRepo,
	}
}

func (s *ProductServices) List(ctx context.Context) ([]*models.Product, error) {
	rows, err := s.productRepo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*models.Product, 0, len(rows))
	for _, p := range rows {
		products = append(products, &models.Product{
			ID:          p.ID,
			BrandID:     p.BrandID,
			CategoryID:  p.CategoryID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			IsActive:    p.IsActive,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return products, nil
}

func (s *ProductServices) DetailByID(ctx context.Context, id string) (*models.Product, error) {
	if id == "" {
		return nil, errors.New(errors.ErrValidation, "invalid id")
	}

	product, err := s.productRepo.GetDetailProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductServices) DetailBySlug(ctx context.Context, slug string) (*models.Product, error) {
	if slug == "" {
		return nil, errors.New(errors.ErrValidation, "invalid slug")
	}

	product, err := s.productRepo.GetDetailProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductServices) Update(ctx context.Context, id, name string) (*models.Product, error) {
	if id == "" {
		return nil, errors.New(errors.ErrValidation, "invalid id")
	}

	if name == "" {
		return nil, errors.New(errors.ErrValidation, "invalid name")
	}

	return nil, nil
}
func (s *ProductServices) Create(ctx context.Context) (*models.Product, error)

func (s *ProductServices) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(errors.ErrValidation, "invalid id")
	}

	if err := s.productRepo.DeleteProduct(ctx, id); err != nil {
		errors.Wrap(errors.ErrInternal, "failed to delete Product", err)
	}

	return nil
}
