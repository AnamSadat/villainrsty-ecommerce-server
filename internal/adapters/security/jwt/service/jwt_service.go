package service

import (
	"time"

	claims "villainrsty-ecommerce-server/internal/adapters/security/jwt/models"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey:          secretKey,
		accessTokenExpiry:  15 * time.Minute,
		refreshTokenExpiry: 7 * 24 * time.Hour,
	}
}

func (s *JWTService) GenerateAccessToken(user *models.User) (string, error) {
	claims := &claims.Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to generate token", err)
	}

	return tokenString, nil
}

func (s *JWTService) GenerateRefreshToken(user *models.User) (string, error) {
	claims := &claims.Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to generate token", err)
	}

	return tokenString, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(errors.ErrUnauthorized, "invalid signin method")
		}

		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorized, "failed to parse token", err)
	}

	claims, ok := token.Claims.(*claims.Claims)
	if !ok || !token.Valid {
		return nil, errors.New(errors.ErrUnauthorized, "invalid token")
	}

	user := &models.User{
		ID:    models.ID(claims.UserID),
		Email: claims.Email,
		Name:  claims.Name,
	}

	return user, nil
}
