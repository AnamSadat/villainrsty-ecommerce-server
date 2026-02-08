package models

import (
	"regexp"

	"villainrsty-ecommerce-server/internal/core/shared/errors"
)

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		User  UserDTO `json:"user"`
		Token string  `json:"token"`
	}

	UserDTO struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	RegisterResponse struct {
		User UserDTO `json:"user"`
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	RefreshTokenResponse struct {
		Token string `json:"token"`
	}
)

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return errors.New(errors.ErrValidation, "email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New(errors.ErrValidation, "email format is invalid")
	}

	if r.Password == "" {
		return errors.New(errors.ErrValidation, "password is required")
	}

	if len(r.Password) < 8 {
		return errors.New(errors.ErrValidation, "password must be at least 8 characters")
	}

	return nil
}

func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return errors.New(errors.ErrValidation, "email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New(errors.ErrValidation, "email is required")
	}

	if r.Password == "" {
		return errors.New(errors.ErrValidation, "password is required")
	}

	if len(r.Password) < 8 {
		return errors.New(errors.ErrValidation, "password must be at least 8 characters")
	}

	if r.Name == "" {
		return errors.New(errors.ErrValidation, "name is required")
	}

	return nil
}

func (r *RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New(errors.ErrValidation, "refresh_token is required")
	}

	return nil
}
