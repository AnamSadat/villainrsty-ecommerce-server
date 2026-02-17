package repository

import (
	"context"
	"errors"

	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/mapper"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"
	appErr "villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
	"villainrsty-ecommerce-server/pkg/validator"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	queries   *sqlc.Queries
	validator *validator.Validator
}

func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{
		queries:   queries,
		validator: validator.NewValidate(),
	}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if err := r.validator.ValidateEmail(email); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invalid email format", err)
	}

	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "user not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get user", err)
	}

	user := mapper.SQLCUserByEmailToDomain(row)

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return nil, appErr.Wrap(appErr.ErrValidation, "invalid id", err)
	}

	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.New(appErr.ErrNotFound, "user not found")
		}

		return nil, appErr.Wrap(appErr.ErrInternal, "failed to get user", err)
	}

	user := mapper.SQLCUserByIDToDomain(row)
	return user, nil
}

func (r *UserRepository) Save(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	exists, err := r.queries.UserExists(ctx, user.Email)
	if err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to check user existence", err)
	}

	if exists {
		return appErr.New(appErr.ErrConflict, "email already registered")
	}

	if err := r.validator.ValidateEmail(user.Email); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid email format", err)
	}

	params := mapper.DomainUserToSQLCParams(user)
	if err := r.queries.CreateUser(ctx, params); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to create user", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid id", err)
	}

	if err := r.queries.DeleteUser(ctx, id); err != nil {
		return appErr.Wrap(appErr.ErrInternal, "failed to delete user", err)
	}

	return nil
}

func (r *UserRepository) Exist(ctx context.Context, email string) (bool, error) {
	if err := r.validator.ValidateEmail(email); err != nil {
		return false, appErr.Wrap(appErr.ErrValidation, "invalid email format", err)
	}

	exists, err := r.queries.UserExists(ctx, email)
	if err != nil {
		return false, appErr.Wrap(appErr.ErrInternal, "failed to check user existence", err)
	}

	return exists, nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, id string, hashed string) error {
	if err := r.validator.ValidateRequired("id", id); err != nil {
		return appErr.Wrap(appErr.ErrValidation, "invalid id", err)
	}

	return nil
}
