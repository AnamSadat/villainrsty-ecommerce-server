package services

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/categories/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/gosimple/slug"
)

type CategoryService struct {
	categoryRepo ports.CategoryRepository
}

func NewCategoriesService(categoryRepo ports.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) List(ctx context.Context) ([]*models.Category, error) {
	rows, err := s.categoryRepo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]*models.Category, 0, len(rows))
	for _, c := range rows {
		categories = append(categories, &models.Category{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	return categories, nil
}

func (s *CategoryService) DetailByID(ctx context.Context, id string) (*models.Category, error) {
	if id == "" {
		return nil, errors.New(errors.ErrValidation, "id is required")
	}

	category, err := s.categoryRepo.GetDetailCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DetailBySlug(ctx context.Context, slug string) (*models.Category, error) {
	if slug == "" {
		return nil, errors.New(errors.ErrValidation, "slug is required")
	}

	category, err := s.categoryRepo.GetDetailCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Update(ctx context.Context, id, name string) (*models.Category, error) {
	if id == "" {
		return nil, errors.New(errors.ErrValidation, "id is required")
	}
	if name == "" {
		return nil, errors.New(errors.ErrValidation, "name is required")
	}

	slugCategory := slug.Make(name)
	exists, err := s.categoryRepo.ExistCategory(ctx, slugCategory)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(errors.ErrConflict, "Category is already")
	}

	if err := s.categoryRepo.UpdateCategory(ctx, id, name, slugCategory); err != nil {
		return nil, err
	}

	return s.categoryRepo.GetDetailCategoryByID(ctx, id)
}

func (s *CategoryService) Create(ctx context.Context, name string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New(errors.ErrValidation, "name is required")
	}

	slugCategory := slug.Make(name)
	category := models.NewCategory(name, slugCategory)
	if err := category.Validate(); err != nil {
		return nil, err
	}

	if err := s.categoryRepo.Save(ctx, category); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to save Category", err)
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(errors.ErrValidation, "id is required")
	}

	if err := s.categoryRepo.DeleteCategory(ctx, id); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete Category", err)
	}

	return nil
}
