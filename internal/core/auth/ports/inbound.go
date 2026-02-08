package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*models.User, string, error)
	Register(ctx context.Context, email, password, name string) (*models.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}
