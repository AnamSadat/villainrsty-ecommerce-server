package models

import (
	"time"

	"villainrsty-ecommerce-server/pkg/validator"
)

type BrandResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=1"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type (
	ListBrandsRequest struct {
		Page  int `query:"page" json:"page"`
		Limit int `query:"limit" json:"limit"`
	}

	ListBrandsResponse struct {
		Brands []BrandResponse `json:"brands"`
		Total  int64           `json:"total"`
	}

	DetailBrandBySlugRequest struct {
		Slug string `param:"slug" validate:"required"`
	}

	DetailBrandBySlugResponse struct {
		Brand []BrandResponse `json:"brand"`
	}
	DetailBrandByIDRequest struct {
		ID string `param:"id" validate:"required"`
	}

	DetailBrandByIDResponse struct {
		Brand []BrandResponse `json:"brand"`
	}

	CreateBrandRequest struct {
		Name string `json:"name" validate:"required,min=1"`
	}

	CreateBrandResponse struct {
		Brand BrandResponse `json:"brand"`
	}

	UpdateBrandRequest struct {
		ID   string `param:"id" validate:"required"`
		Name string `json:"name" validate:"required,min=1"`
	}

	UpdateBrandResponse struct {
		Brand BrandResponse `json:"brand"`
	}

	DeleteBrandRequest struct {
		ID string `param:"id" validation:"required"`
	}
)

func (r *ListBrandsRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *DetailBrandByIDRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *DetailBrandBySlugRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *CreateBrandRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *UpdateBrandRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *DeleteBrandRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}
