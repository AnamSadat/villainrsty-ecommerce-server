package models

import (
	"time"

	"villainrsty-ecommerce-server/pkg/validator"
)

type Category struct {
	ID        ID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCategory(name, slug string) *Category {
	now := time.Now()
	return &Category{
		ID:        NewID(),
		Name:      name,
		Slug:      slug,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *Category) Validate() error {
	v := validator.NewValidate()
	if err := v.ValidateName(p.Name); err != nil {
		return err
	}

	return nil
}
