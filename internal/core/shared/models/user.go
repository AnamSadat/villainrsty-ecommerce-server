package models

import (
	"time"

	"villainrsty-ecommerce-server/pkg/validator"
)

type User struct {
	ID        ID
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, password, name string) *User {
	now := time.Now()
	return &User{
		ID:        NewID(),
		Email:     email,
		Password:  password,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Validate() error {
	v := validator.NewValidate()

	if err := v.ValidateEmail(u.Email); err != nil {
		return err
	}

	if err := v.ValidatePassword(u.Password); err != nil {
		return err
	}

	if err := v.ValidateName(u.Name); err != nil {
		return err
	}

	if err := u.ID.Validate(); err != nil {
		return err
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
