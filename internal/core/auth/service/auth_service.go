package service

import (
	"context"

	"villainrsty-ecommerce-server/internal/core/auth/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type AuthService struct {
	userRepo   ports.UserRepository
	hasher     ports.PasswordHasher
	jwtService ports.JWTService
}

func NewAuthService(
	userRepo ports.UserRepository,
	hasher ports.PasswordHasher,
	jwtService ports.JWTService,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		hasher:     hasher,
		jwtService: jwtService,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	if email == "" || password == "" {
		return nil, "", errors.New(errors.ErrValidation, "email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	if !s.hasher.Verify(user.Password, password) {
		return nil, "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	token, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, "", errors.New(errors.ErrInternal, "failed to generate token")
	}

	return user, token, nil
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

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", errors.New(errors.ErrValidation, "refresh token is required")
	}

	user, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New(errors.ErrUnauthorized, "invalid refresh token")
	}

	token, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to generate refresh token", err)
	}

	return token, nil
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
