package app

import (
	"log/slog"

	"villainrsty-ecommerce-server/internal/adapters/http/auth/handler"
	"villainrsty-ecommerce-server/internal/adapters/notifications/smtp"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/repository"
	tokenHasher "villainrsty-ecommerce-server/internal/adapters/security/hasher"
	jwtService "villainrsty-ecommerce-server/internal/adapters/security/jwt/service"
	"villainrsty-ecommerce-server/internal/adapters/security/password"
	"villainrsty-ecommerce-server/internal/config"
	"villainrsty-ecommerce-server/internal/core/auth/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	AuthHandler *handler.AuthHandler
}

func New(cfg config.Config, db *pgxpool.Pool, logger *slog.Logger) *Container {
	queries := postgres.NewQueries(db)
	userRepo := repository.NewUserRepository(queries)
	refreshTokenRepo := repository.NewRefreshTokenRepository(queries)
	passwordResetRepo := repository.NewPasswordResetTokenRepository(queries)
	hasher := password.NewBcryptHasher()
	tokenHasher := tokenHasher.NewSHA256TokenHasher()
	jwtService := jwtService.NewJWTService(cfg.CookieSecret)
	emailSender := smtp.NewEmailSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFromEmail,
		cfg.SMTPFromName)
	authService := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
		passwordResetRepo,
		emailSender,
		hasher,
		tokenHasher,
		jwtService,
		logger,
		cfg.ResetPasswordURL,
		cfg.ResetPasswordTTL,
	)
	authHandler := handler.NewAuthHandler(authService, logger)

	return &Container{
		AuthHandler: authHandler,
	}
}
