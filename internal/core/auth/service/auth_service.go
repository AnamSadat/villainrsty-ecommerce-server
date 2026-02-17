package service

import (
	"context"
	"log/slog"
	"time"

	"villainrsty-ecommerce-server/internal/core/auth/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type AuthService struct {
	userRepo         ports.UserRepository
	refreshTokenRepo ports.RefreshTokenRepository
	hasher           ports.PasswordHasher
	tokenHasher      ports.TokenHasher
	jwtService       ports.JWTService
	logger           *slog.Logger
}

func NewAuthService(
	userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository,
	hasher ports.PasswordHasher,
	tokenHasher ports.TokenHasher,
	jwtService ports.JWTService,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		hasher:           hasher,
		tokenHasher:      tokenHasher,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		logger:           logger,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string, rememberMe bool) (*models.User, string, string, error) {
	if email == "" || password == "" {
		return nil, "", "", errors.New(errors.ErrValidation, "email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	if user == nil {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	if !s.hasher.Verify(user.Password, password) {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, "", "", errors.New(errors.ErrInternal, "failed to generate access token")
	}

	refreshTokenString, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", errors.New(errors.ErrInternal, "failed to generate refresh token")
	}

	tokenHash, err := s.tokenHasher.Hash(refreshTokenString)
	s.logger.Info("refresh token debug (login)",
		"refresh_len", len(refreshTokenString),
		"hash_len", len(tokenHash),
		"hash_prefix", tokenHash[:8],
	)

	if err != nil {
		s.logger.Warn("token hash", "error", err)
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to hash token", err)
	}

	ttl := 24 * time.Hour
	if rememberMe {
		ttl = 30 * 24 * time.Hour
		s.logger.Info("ttl", "ttl sekarang: ", ttl)
	}

	refreshToken := models.NewRefreshToken(user.ID, tokenHash, ttl)
	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to save refresh token", err)
	}

	return user, accessToken, refreshTokenString, nil
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, errors.New(errors.ErrValidation, "email and password are required")
	}

	exists, err := s.userRepo.Exist(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(errors.ErrConflict, "email already registered")
	}

	user := models.NewUser(email, password, name)
	if !user.IsPasswordValid(password) {
		return nil, errors.New(errors.ErrValidation, "password must contain uppercase, lowercase and number")
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password", err)
	}

	user.Password = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New(errors.ErrValidation, "refresh token is required")
	}

	tokenHash, err := s.tokenHasher.Hash(refreshToken)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to hash refresh token", err)
	}

	dbToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return errors.New(errors.ErrUnauthorized, "refresh token not found")
	}

	if err := s.refreshTokenRepo.Revoke(ctx, dbToken.ID); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to revoke token", err)
	}

	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", errors.New(errors.ErrValidation, "refresh token is required")
	}

	user, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New(errors.ErrUnauthorized, "invalid refresh token")
	}

	tokenHash, err := s.tokenHasher.Hash(refreshToken)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to hash refresh token", err)
	}
	s.logger.Info("refresh token debug (incoming)",
		"incoming_len", len(refreshToken),
		"hash_len", len(tokenHash),
		"hash_prefix", tokenHash[:8],
	)

	dbToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return "", "", errors.New(errors.ErrUnauthorized, "refresh token not found")
	}

	if !dbToken.IsValid() {
		return "", "", errors.New(errors.ErrUnauthorized, "refresh token is expired or revoked")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to generate access token", err)
	}

	newRefreshTokenString, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to generate refresh token", err)
	}

	newTokenHash, err := s.tokenHasher.Hash(newRefreshTokenString)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to hash refresh token", err)
	}

	newRefreshToken := models.NewRefreshToken(user.ID, newTokenHash, 7*24*time.Hour)
	if err := s.refreshTokenRepo.Save(ctx, newRefreshToken); err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to save refresh token", err)
	}

	if err := s.refreshTokenRepo.Revoke(ctx, dbToken.ID); err != nil {
		s.logger.Error("failed to revoke old refresh token", err)
	}

	return accessToken, newRefreshTokenString, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	if token == "" {
		return nil, errors.New(errors.ErrValidation, "token is required")
	}

	user, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, errors.New(errors.ErrUnauthorized, "invalid token")
	}

	return user, nil
}

func (s *AuthService) UpdateUserPassword()
