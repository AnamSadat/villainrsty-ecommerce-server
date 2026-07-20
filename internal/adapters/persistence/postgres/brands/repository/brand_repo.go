package repository

import (
	"context"
	"errors"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/brands/mapper"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	appErr "villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
	"villainrsty-ecommerce-server/pkg/validator"

	"github.com/jackc/pgx/v5"
)

type BrandRepository struct {
	queries   *sqlc.Queries
	validator *validator.Validator
}

func NewBrandRepository(queries *sqlc.Queries) *BrandRepository {
	return &BrandRepository{
		queries:   queries,
		validator: validator.NewValidate(),
	}
}

func (r *BrandRepository) Save(ctx context.Context, brand *models.Brand) error {
	if err := brand.Validate(); err != nil {
		return err
	}

	exists, err := r.queries.BrandExists(ctx, brand.Slug)
	if err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to check brand exists", err)
	}

	if exists {
		return appErr.New(appErr.ErrConflict, "brand is already")
	}

	params := mapper.DomainBrandToSQLParams(brand)
	if err := r.queries.CreateBrand(ctx, params); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to create brand", err)
	}

	return nil
}

func (r *BrandRepository) GetBrands(ctx context.Context) ([]*models.Brand, error) {
	rows, err := r.queries.GetAllBrands(ctx)
	if err != nil {
		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get all brands", err)
	}

	brands := make([]*models.Brand, len(rows))
	for i, row := range rows {
		brands[i] = mapper.SQLBrandBySlugToDomain(row)
	}

	return brands, nil
}

func (r *BrandRepository) GetDetailBrandByID(ctx context.Context, id string) (*models.Brand, error) {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invaid id", err)
	}

	row, err := r.queries.GetBrandByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "brand not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get brand", err)
	}

	brand := mapper.SQLBrandBySlugToDomain(row)
	return brand, nil
}

func (r *BrandRepository) GetDetailBrandBySlug(ctx context.Context, slug string) (*models.Brand, error) {
	if err := r.validator.ValidateRequired("slug", slug); err != nil {
		return nil, appErr.Wrap(appErr.ErrInternal, "invalid slug", err)
	}

	row, err := r.queries.GetBrandBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "brand slug not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get brand slug", err)
	}

	brandSlug := mapper.SQLBrandBySlugToDomain(row)

	return brandSlug, nil
}

func (r *BrandRepository) UpdateBrand(ctx context.Context, id, name, slug string) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.New(appErr.ErrValidation, "invalid id")
	}

	if err := r.validator.ValidateName(name); err != nil {
		return appErr.New(appErr.ErrValidation, "invalid name")
	}

	err := r.queries.UpdateBrandByID(ctx, sqlc.UpdateBrandByIDParams{
		ID:   id,
		Name: name,
		Slug: slug,
	})
	if err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to update brand name", err)
	}

	return nil
}

func (r *BrandRepository) DeleteBrand(ctx context.Context, id string) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid id", err)
	}

	if err := r.queries.DeleteBrandByID(ctx, id); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to delete brand", err)
	}

	return nil
}

func (r *BrandRepository) ExistBrand(ctx context.Context, slug string) (bool, error) {
	if err := r.validator.ValidateRequired("slug", slug); err != nil {
		return false, appErr.Wrap(appErr.ErrValidation, "invalid slug", err)
	}

	exists, err := r.queries.BrandExists(ctx, slug)
	if err != nil {
		return false, appErr.Wrap(appErr.ErrInternal, "failed to check brand existence", err)
	}

	return exists, nil
}
