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
		UpdateUserPassword(ctx context.Context, id models.ID, hashed string) error
	}

	RefreshTokenRepository interface {
		Save(ctx context.Context, token *models.RefreshToken) error
		GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
		GetByUserID(ctx context.Context, userID models.ID) ([]*models.RefreshToken, error)
		Revoke(ctx context.Context, tokenID models.ID) error
		DeleteExpired(ctx context.Context) error
	}

	PasswordHasher interface {
		Hash(password string) (string, error)
		Verify(hash, password string) bool
	}

	TokenHasher interface {
		Hash(token string) (string, error)
	}

	JWTService interface {
		GenerateAccessToken(user *models.User) (string, error)
		GenerateRefreshToken(user *models.User) (string, error)
		ValidateToken(token string) (*models.User, error)
	}

	PasswordResetTokenRepository interface {
		Save(ctx context.Context, t *models.PasswordResetToken) error
		GetByTokenHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error)
		MarkUsed(ctx context.Context, id models.ID) error
	}
)
