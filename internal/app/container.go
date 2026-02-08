package app

import (
	"villainrsty-ecommerce-server/internal/adapters/http/auth/handler"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres"
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/auth/repository"
	jwtService "villainrsty-ecommerce-server/internal/adapters/security/jwt/service"
	"villainrsty-ecommerce-server/internal/adapters/security/password"
	"villainrsty-ecommerce-server/internal/config"
	"villainrsty-ecommerce-server/internal/core/auth/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	AuthHandler *handler.AuthHandler
	AuthService *service.AuthService
}

func New(cfg config.Config, db *pgxpool.Pool) *Container {
	queries := postgres.NewQueries(db)
	userRepo := repository.NewUserRepository(queries)
	hasher := password.NewBcryptHasher()
	jwtService := jwtService.NewJWTService(cfg.CookieSecret)
	authService := service.NewAuthService(userRepo, hasher, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	return &Container{
		AuthHandler: authHandler,
		AuthService: authService,
	}
}
