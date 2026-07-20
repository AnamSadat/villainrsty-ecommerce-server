package repository

import (
	"context"
	"errors"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/categories/mapper"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	appErr "villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
	"villainrsty-ecommerce-server/pkg/validator"

	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	queries   *sqlc.Queries
	validator *validator.Validator
}

func NewCategoryRepository(queries *sqlc.Queries) *CategoryRepository {
	return &CategoryRepository{
		queries:   queries,
		validator: validator.NewValidate(),
	}
}

func (r *CategoryRepository) Save(ctx context.Context, category *models.Category) error {
	if err := category.Validate(); err != nil {
		return err
	}

	exists, err := r.queries.CategoryExists(ctx, category.Slug)
	if err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to check category exists", err)
	}

	if exists {
		return appErr.New(appErr.ErrConflict, "category is already")
	}

	params := mapper.DomainCategoryToSQLParams(category)
	if err := r.queries.CreateCategory(ctx, params); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to create category", err)
	}

	return nil
}

func (r *CategoryRepository) GetCategories(ctx context.Context) ([]*models.Category, error) {
	rows, err := r.queries.GetAllCategories(ctx)
	if err != nil {
		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get all categories", err)
	}

	categories := make([]*models.Category, len(rows))
	for i, row := range rows {
		categories[i] = mapper.SQLCategoryBySlugToDomain(row)
	}

	return categories, nil
}

func (r *CategoryRepository) GetDetailCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invaid id", err)
	}

	row, err := r.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "category not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get category", err)
	}

	category := mapper.SQLCategoryBySlugToDomain(row)
	return category, nil
}

func (r *CategoryRepository) GetDetailCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	if err := r.validator.ValidateRequired("slug", slug); err != nil {
		return nil, appErr.Wrap(appErr.ErrInternal, "invalid slug", err)
	}

	row, err := r.queries.GetCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "category slug not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get category slug", err)
	}

	categorySlug := mapper.SQLCategoryBySlugToDomain(row)

	return categorySlug, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, id, name, slug string) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.New(appErr.ErrValidation, "invalid id")
	}

	if err := r.validator.ValidateName(name); err != nil {
		return appErr.New(appErr.ErrValidation, "invalid name")
	}

	err := r.queries.UpdateCategoryByID(ctx, sqlc.UpdateCategoryByIDParams{
		ID:   id,
		Name: name,
		Slug: slug,
	})
	if err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to update category name", err)
	}

	return nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid id", err)
	}

	if err := r.queries.DeleteCategoryByID(ctx, id); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to delete category", err)
	}

	return nil
}

func (r *CategoryRepository) ExistCategory(ctx context.Context, slug string) (bool, error) {
	if err := r.validator.ValidateRequired("slug", slug); err != nil {
		return false, appErr.Wrap(appErr.ErrValidation, "invalid slug", err)
	}

	exists, err := r.queries.CategoryExists(ctx, slug)
	if err != nil {
		return false, appErr.Wrap(appErr.ErrInternal, "failed to check category existence", err)
	}

	return exists, nil
}
