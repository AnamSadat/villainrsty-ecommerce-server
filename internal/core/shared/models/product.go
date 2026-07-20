package models

import (
	"time"

	"villainrsty-ecommerce-server/pkg/validator"
)

type Product struct {
	ID          ID
	BrandID     ID
	CategoryID  ID
	Name        string
	Slug        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(brandID, categoryID ID, name, slug, description string, isActive bool) *Product {
	now := time.Now()
	return &Product{
		ID:          NewID(),
		BrandID:     brandID,
		CategoryID:  categoryID,
		Name:        name,
		Slug:        slug,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Product) Validate() error {
	v := validator.NewValidate()

	if err := p.ID.Validate(); err != nil {
		return err
	}

	if err := p.BrandID.Validate(); err != nil {
		return err
	}

	if err := p.CategoryID.Validate(); err != nil {
		return err
	}

	if err := v.ValidateName(p.Name); err != nil {
		return err
	}

	if p.Slug == "" {
		return validator.NewValidate().ValidateRequired("slug", p.Slug)
	}

	if p.Description == "" {
		return validator.NewValidate().ValidateRequired("description", p.Description)
	}

	return nil
}
