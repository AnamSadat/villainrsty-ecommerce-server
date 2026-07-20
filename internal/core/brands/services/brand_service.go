package services

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/brands/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/gosimple/slug"
)

type BrandService struct {
	brandRepo ports.BrandRepository
}

func NewBrandService(brandRepo ports.BrandRepository) *BrandService {
	return &BrandService{brandRepo: brandRepo}
}

func (s *BrandService) List(ctx context.Context) ([]*models.Brand, error) {
	rows, err := s.brandRepo.GetBrands(ctx)
	if err != nil {
		return nil, err
	}

	brands := make([]*models.Brand, 0, len(rows))
	for _, b := range rows {
		brands = append(brands, &models.Brand{
			ID:        b.ID,
			Name:      b.Name,
			Slug:      b.Slug,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	return brands, nil
}

func (s *BrandService) DetailByID(ctx context.Context, id string) (*models.Brand, error) {
	if id == "" {
		return nil, errors.New(errors.ErrValidation, "id is required")
	}

	brand, err := s.brandRepo.GetDetailBrandByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return brand, nil
}

func (s *BrandService) DetailBySlug(ctx context.Context, slug string) (*models.Brand, error) {
	if slug == "" {
		return nil, errors.New(errors.ErrValidation, "slug is required")
	}

	brand, err := s.brandRepo.GetDetailBrandBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return brand, nil
}

func (s *BrandService) Update(ctx context.Context, id, name string) (*models.Brand, error) {
	if id == "" {
		return nil, errors.New(errors.ErrValidation, "id is required")
	}
	if name == "" {
		return nil, errors.New(errors.ErrValidation, "name is required")
	}

	exists, err := s.brandRepo.ExistBrand(ctx, name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(errors.ErrConflict, "Brand is already")
	}

	slugBrand := slug.Make(name)
	if err := s.brandRepo.UpdateBrand(ctx, id, name, slugBrand); err != nil {
		return nil, err
	}

	return s.brandRepo.GetDetailBrandByID(ctx, id)
}

func (s *BrandService) Create(ctx context.Context, name string) (*models.Brand, error) {
	if name == "" {
		return nil, errors.New(errors.ErrValidation, "name is required")
	}

	slugBrand := slug.Make(name)
	brand := models.NewBrand(name, slugBrand)

	if err := brand.Validate(); err != nil {
		return nil, err
	}

	if err := s.brandRepo.Save(ctx, brand); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to save brand", err)
	}

	return brand, nil
}

func (s *BrandService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(errors.ErrValidation, "id is required")
	}

	if err := s.brandRepo.DeleteBrand(ctx, id); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete brand", err)
	}

	return nil
}
