package models

import (
	"time"

	"villainrsty-ecommerce-server/pkg/validator"
)

type Brand struct {
	ID        ID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBrand(name, slug string) *Brand {
	now := time.Now()
	return &Brand{
		ID:        NewID(),
		Name:      name,
		Slug:      slug,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *Brand) Validate() error {
	v := validator.NewValidate()
	if err := v.ValidateName(b.Name); err != nil {
		return err
	}

	if err := b.ID.Validate(); err != nil {
		return err
	}

	return nil
}
