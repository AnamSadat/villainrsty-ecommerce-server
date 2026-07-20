package models

import (
	"time"

	"villainrsty-ecommerce-server/pkg/validator"
)

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=1"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type (
	ListCategoriesRequest struct {
		Page  int `query:"page" json:"page"`
		Limit int `query:"limit" json:"limit"`
	}

	ListCategoriesResponse struct {
		Categories []CategoryResponse `json:"categories"`
		Total      int64              `json:"total"`
	}

	DetailCategoryBySlugRequest struct {
		Slug string `param:"slug" validate:"required"`
	}

	DetailCategoryBySlugResponse struct {
		Category []CategoryResponse `json:"category"`
	}

	DetailCategoryByIDRequest struct {
		ID string `param:"id" validate:"required"`
	}

	DetailCategoryByIDResponse struct {
		Category []CategoryResponse `json:"category"`
	}

	CreateCategoryRequest struct {
		Name string `json:"name" validate:"required,min=1"`
	}

	CreateCategoryResponse struct {
		Category CategoryResponse `json:"category"`
	}

	UpdateCategoryRequest struct {
		ID   string `param:"id" validate:"required"`
		Name string `json:"name" validate:"required,min=1"`
	}

	UpdateCategoryResponse struct {
		Category CategoryResponse `json:"category"`
	}

	DeleteCategoryRequest struct {
		ID string `param:"id" validation:"required"`
	}
)

func (r *ListCategoriesRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *DetailCategoryByIDRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *DetailCategoryBySlugRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *CreateCategoryRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *UpdateCategoryRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *DeleteCategoryRequest) Validate() error {
	v := validator.NewValidate()
	return v.ValidateStruct(r)
}
