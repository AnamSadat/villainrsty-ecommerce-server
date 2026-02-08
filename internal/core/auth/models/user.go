package models

import (
	"time"

	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type User struct {
	ID        models.ID
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, password, name string) *User {
	now := time.Now()
	return &User{
		ID:        models.NewID(),
		Email:     email,
		Password:  password,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New(errors.ErrValidation, "email is required")
	}

	if u.Name == "" {
		return errors.New(errors.ErrValidation, "name is required")
	}

	if u.Password == "" {
		return errors.New(errors.ErrValidation, "password is required")
	}

	return nil
}

func (u *User) IsPasswordValid(plainPassword string) bool {
	if len(plainPassword) < 8 {
		return false
	}

	hasUpper := false
	for _, c := range plainPassword {
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
			break
		}
	}

	if !hasUpper {
		return false
	}

	hasLower := false
	for _, c := range plainPassword {
		if c >= 'a' && c <= 'z' {
			hasLower = true
			break
		}
	}

	if !hasLower {
		return false
	}

	hasNumber := false
	for _, c := range plainPassword {
		if c >= '0' && c <= '9' {
			hasNumber = true
			break
		}
	}

	if !hasNumber {
		return false
	}

	return true
}
