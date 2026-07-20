package app

import (
	"log/slog"

	"villainrsty-ecommerce-server/internal/adapters/notifications/smtp"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres"
	authRepo "villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/repository"
	brandRepo "villainrsty-ecommerce-server/internal/adapters/persistence/postgres/brands/repository"
	categoryRepo "villainrsty-ecommerce-server/internal/adapters/persistence/postgres/categories/repository"
	"villainrsty-ecommerce-server/internal/adapters/security/password"
	"villainrsty-ecommerce-server/internal/config"
	"villainrsty-ecommerce-server/internal/core/auth/ports"
	authService "villainrsty-ecommerce-server/internal/core/auth/services"
	brandService "villainrsty-ecommerce-server/internal/core/brands/services"
	categoryService "villainrsty-ecommerce-server/internal/core/categories/services"

	authHandler "villainrsty-ecommerce-server/internal/adapters/http/auth/handler"
	brandHandler "villainrsty-ecommerce-server/internal/adapters/http/brands/handler"
	categoryHandler "villainrsty-ecommerce-server/internal/adapters/http/categories/handler"
	rbacHandler "villainrsty-ecommerce-server/internal/adapters/http/rbac/handler"
	authorization "villainrsty-ecommerce-server/internal/adapters/security/authorization"
	tokenHasher "villainrsty-ecommerce-server/internal/adapters/security/hasher"
	jwtService "villainrsty-ecommerce-server/internal/adapters/security/jwt/service"

	"github.com/casbin/casbin/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	AuthHandler     *authHandler.AuthHandler
	RbacHandler     *rbacHandler.RbacHandler
	BrandHandler    *brandHandler.BrandHandler
	CategoryHandler *categoryHandler.CategoryHandler
	AuthService     ports.AuthService
	JWTService      ports.JWTService
	CasbinEnforcer  *casbin.Enforcer
	Logger          *slog.Logger
}

func New(cfg config.Config, db *pgxpool.Pool, logger *slog.Logger) *Container {
	queries := postgres.NewQueries(db)
	userRepo := authRepo.NewUserRepository(queries)
	brandRepo := brandRepo.NewBrandRepository(queries)
	categoryRepo := categoryRepo.NewCategoryRepository(queries)
	refreshTokenRepo := authRepo.NewRefreshTokenRepository(queries)
	passwordResetRepo := authRepo.NewPasswordResetTokenRepository(queries)
	twoFactorOTPRepo := authRepo.NewTwoFactorOTPRepository(queries)
	userRoleRepo := authRepo.NewUserRoleRepository(queries)
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

	authService := authService.NewAuthService(
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

	brandService := brandService.NewBrandService(brandRepo)
	categoryService := categoryService.NewCategoriesService(categoryRepo)

	authHTTPHandler := authHandler.NewAuthHandler(authService, logger)
	rbacHTTPHandler := rbacHandler.NewRbacHandler()
	brandHTTPHandler := brandHandler.NewBrandHandler(brandService, logger)
	categoryHTTPHandler := categoryHandler.NewCategoryHandler(categoryService, logger)

	return &Container{
		AuthHandler:     authHTTPHandler,
		RbacHandler:     rbacHTTPHandler,
		BrandHandler:    brandHTTPHandler,
		CategoryHandler: categoryHTTPHandler,
		AuthService:     authService,
		JWTService:      jwtService,
		CasbinEnforcer:  enforcer,
		Logger:          logger,
	}
}
