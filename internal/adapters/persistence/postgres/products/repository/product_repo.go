package repository

import (
	"context"
	"errors"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/products/mapper"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	"villainrsty-ecommerce-server/internal/core/shared/models"
	"villainrsty-ecommerce-server/pkg/validator"

	appErr "villainrsty-ecommerce-server/internal/core/shared/errors"

	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	queries   *sqlc.Queries
	validator *validator.Validator
}

func NewProductRepository(query *sqlc.Queries) *ProductRepository {
	return &ProductRepository{queries: query, validator: validator.NewValidate()}
}

func (r *ProductRepository) Save(ctx context.Context, product *models.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}

	exists, err := r.queries.ProductExists(ctx, product.Slug)
	if err != nil {
		return appErr.Wrap(appErr.ErrInternal, "Failed to check product exist", err)
	}

	if exists {
		return appErr.New(appErr.ErrConflict, "Product is already")
	}

	params := mapper.DomainProductToSQLParams(product)
	if err := r.queries.CreateProduct(ctx, params); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "Failed to add product", err)
	}

	return nil
}

func (r *ProductRepository) GetProducts(ctx context.Context) ([]*models.Product, error) {
	rows, err := r.queries.GetAllProducts(ctx)
	if err != nil {
		return nil, appErr.Wrap(appErr.ErrInternal, "Failed to get products", err)
	}

	products := make([]*models.Product, len(rows))
	for i, row := range rows {
		products[i] = mapper.SQLProductToDamain(row)
	}

	return products, nil
}

func (r *ProductRepository) GetDetailProductByID(ctx context.Context, id string) (*models.Product, error) {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invalid id", err)
	}

	row, err := r.queries.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "Product not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get product by id", err)
	}

	product := mapper.SQLProductToDamain(row)
	return product, nil
}

func (r *ProductRepository) GetDetailProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	if err := r.validator.ValidateRequired("slug", slug); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invalid slug", err)
	}

	row, err := r.queries.GetProductBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "Product not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get product by slug", err)
	}

	product := mapper.SQLProductToDamain(row)
	return product, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id, brandID, categoryID, name, slug, description string, isActive bool) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.New(appErr.ErrValidation, "invalid id")
	}

	if err := r.validator.ValidateName(name); err != nil {
		return appErr.New(appErr.ErrValidation, "invalid name")
	}

	err := r.queries.UpdateProductByID(ctx, sqlc.UpdateProductByIDParams{
		ID:          id,
		BrandID:     brandID,
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		IsActive:    isActive,
	})
	if err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to update product", err)
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.New(appErr.ErrValidation, "invalid id")
	}

	if err := r.queries.DeleteProductByID(ctx, id); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to delete product by id", err)
	}

	return nil
}

func (r *ProductRepository) ExistProduct(ctx context.Context, slug string) (bool, error) {
	if err := r.validator.ValidateRequired("slug", slug); err != nil {
		return false, appErr.New(appErr.ErrValidation, "invalid slug")
	}

	exists, err := r.queries.ProductExists(ctx, slug)
	if err != nil {
		return false, appErr.Wrap(appErr.ErrInternal, "failed to check product existence", err)
	}

	return exists, nil
}
