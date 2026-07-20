package handler

import (
	"log/slog"
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/categories/models"
	"villainrsty-ecommerce-server/internal/adapters/http/lib/httpx"
	"villainrsty-ecommerce-server/internal/core/categories/ports"

	"villainrsty-ecommerce-server/internal/core/shared/errors"
	sharedModel "villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	categoryService ports.CategoryService
	logger          *slog.Logger
}

func NewCategoryHandler(service ports.CategoryService, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{categoryService: service, logger: logger}
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.List(r.Context())
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	dtos := make([]models.CategoryResponse, len(categories))
	for i, cat := range categories {
		dtos[i] = mapCategoryToDTO(cat)
	}

	httpx.Success(w, http.StatusOK, "Get all category success", dtos)
}

func (h *CategoryHandler) GetDetailCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		httpx.HandleError(w, errors.New(errors.ErrValidation, "slug is required"), h.logger)
	}

	category, err := h.categoryService.DetailBySlug(r.Context(), slug)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	dtos := mapCategoryToDTO(category)

	httpx.Success(w, http.StatusOK, "Get detail category by slug success", dtos)
}

func (h *CategoryHandler) GetDetailCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.HandleError(w, errors.New(errors.ErrValidation, "id is required"), h.logger)
	}

	category, err := h.categoryService.DetailByID(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}
	dtos := mapCategoryToDTO(category)

	httpx.Success(w, http.StatusOK, "Get detail category by slug success", dtos)
}

func (h *CategoryHandler) AddCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	category, err := h.categoryService.Create(r.Context(), req.Name)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	resp := models.CreateCategoryResponse{Category: mapCategoryToDTO(category)}

	httpx.Success(w, http.StatusOK, "Add category success", resp)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateCategoryRequest
	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.HandleError(w, errors.New(errors.ErrValidation, "id is required"), h.logger)
	}

	req.ID = id

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	category, err := h.categoryService.Update(r.Context(), req.ID, req.Name)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	httpx.Success(w, http.StatusOK, "success update category", category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.HandleError(w, errors.New(errors.ErrValidation, "id is required"), h.logger)
	}

	err := h.categoryService.Delete(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	httpx.Success(w, http.StatusOK, "Delete category success", "")
}

func mapCategoryToDTO(category *sharedModel.Category) models.CategoryResponse {
	return models.CategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
