package ports

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type (
	AuthService interface {
		Login(ctx context.Context, email, password string, rememberMe bool) (*models.User, string, string, error)
		Register(ctx context.Context, email, password, name string) (*models.User, error)
		RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
		ValidateToken(ctx context.Context, token string) (*models.User, error)
		Logout(ctx context.Context, refreshToken string) error
		RequestPasswordReset(ctx context.Context, email string) error
		ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
	}

	EmailSender interface {
		SendPasswordReset(ctx context.Context, toEmail, resetLink string) error
	}
)
