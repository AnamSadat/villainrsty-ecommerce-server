package app

import (
	"log/slog"

	"villainrsty-ecommerce-server/internal/adapters/notifications/smtp"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/repository"
	"villainrsty-ecommerce-server/internal/adapters/security/password"
	"villainrsty-ecommerce-server/internal/config"
	"villainrsty-ecommerce-server/internal/core/auth/ports"
	"villainrsty-ecommerce-server/internal/core/auth/service"

	authHandler "villainrsty-ecommerce-server/internal/adapters/http/auth/handler"
	rbacHandler "villainrsty-ecommerce-server/internal/adapters/http/rbac/handler"
	authorization "villainrsty-ecommerce-server/internal/adapters/security/authorization"
	tokenHasher "villainrsty-ecommerce-server/internal/adapters/security/hasher"
	jwtService "villainrsty-ecommerce-server/internal/adapters/security/jwt/service"

	"github.com/casbin/casbin/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	AuthHandler    *authHandler.AuthHandler
	RbacHandler    *rbacHandler.RbacHandler
	AuthService    ports.AuthService
	JWTService     ports.JWTService
	CasbinEnforcer *casbin.Enforcer
}

func New(cfg config.Config, db *pgxpool.Pool, logger *slog.Logger) *Container {
	queries := postgres.NewQueries(db)
	userRepo := repository.NewUserRepository(queries)
	refreshTokenRepo := repository.NewRefreshTokenRepository(queries)
	passwordResetRepo := repository.NewPasswordResetTokenRepository(queries)
	twoFactorOTPRepo := repository.NewTwoFactorOTPRepository(queries)
	userRoleRepo := repository.NewUserRoleRepository(queries)
	hasher := password.NewBcryptHasher()
	tokenHasher := tokenHasher.NewSHA256TokenHasher()
	jwtService := jwtService.NewJWTService(cfg.CookieSecret)
	enforcer, err := authorization.NewEnforcer(cfg.CasbinModelPath, cfg.CasbinPolicyPath)
	if err != nil {
		panic(err)
	}

	emailSender := smtp.NewEmailSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFromEmail,
		cfg.SMTPFromName,
	)

	authService := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
		passwordResetRepo,
		twoFactorOTPRepo,
		userRoleRepo,
		emailSender,
		hasher,
		tokenHasher,
		jwtService,
		logger,
		cfg.ResetPasswordURL,
		cfg.ResetPasswordTTL,
		cfg.TwoFactorOTPTTL,
	)

	authHTTPHandler := authHandler.NewAuthHandler(authService, logger)
	rbacHTTPHandler := rbacHandler.NewRbacHandler()

	return &Container{
		AuthHandler:    authHTTPHandler,
		RbacHandler:    rbacHTTPHandler,
		AuthService:    authService,
		JWTService:     jwtService,
		CasbinEnforcer: enforcer,
	}
}
