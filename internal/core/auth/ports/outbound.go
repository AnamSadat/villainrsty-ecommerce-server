package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	UserRepository interface {
		GetByEmail(ctx context.Context, email string) (*models.User, error)
		GetByID(ctx context.Context, id string) (*models.User, error)
		Save(ctx context.Context, user *models.User) error
		Delete(ctx context.Context, id string) error
		Exist(ctx context.Context, email string) (bool, error)
	}

	PasswordHasher interface {
		Hash(password string) (string, error)
		Verify(hash, password string) bool
	}

	JWTService interface {
		GenerateAccessToken(user *models.User) (string, error)
		GenerateRefreshToken(user *models.User) (string, error)
		ValidateToken(token string) (*models.User, error)
	}
)
