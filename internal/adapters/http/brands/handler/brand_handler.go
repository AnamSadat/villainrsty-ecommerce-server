package handler

import (
	"log/slog"
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/brands/models"
	"villainrsty-ecommerce-server/internal/adapters/http/lib/httpx"
	"villainrsty-ecommerce-server/internal/core/brands/ports"

	"villainrsty-ecommerce-server/internal/core/shared/errors"
	sharedModel "villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/go-chi/chi/v5"
)

type BrandHandler struct {
	brandService ports.BrandService
	logger       *slog.Logger
}

func NewBrandHandler(service ports.BrandService, logger *slog.Logger) *BrandHandler {
	return &BrandHandler{brandService: service, logger: logger}
}

func (h *BrandHandler) GetAllBrand(w http.ResponseWriter, r *http.Request) {
	brands, err := h.brandService.List(r.Context())
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	dtos := make([]models.BrandResponse, len(brands))
	for i, cat := range brands {
		dtos[i] = mapBrandToDTO(cat)
	}

	httpx.Success(w, http.StatusOK, "Get all brand success", dtos)
}

func (h *BrandHandler) GetDetailBrandBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		httpx.HandleError(w, errors.New(errors.ErrValidation, "slug is required"), h.logger)
		return
	}

	brand, err := h.brandService.DetailBySlug(r.Context(), slug)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	dtos := mapBrandToDTO(brand)

	httpx.Success(w, http.StatusOK, "Get detail brand by slug success", dtos)
}

func (h *BrandHandler) GetDetailBrandByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.HandleError(w, errors.New(errors.ErrValidation, "id is required"), h.logger)
	}

	brand, err := h.brandService.DetailByID(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	dtos := mapBrandToDTO(brand)

	httpx.Success(w, http.StatusOK, "Get detail brand by slug success", dtos)
}

func (h *BrandHandler) AddBrand(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBrandRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	brand, err := h.brandService.Create(r.Context(), req.Name)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	resp := models.CreateBrandResponse{Brand: mapBrandToDTO(brand)}

	httpx.Success(w, http.StatusOK, "Add brand success", resp)
}

func (h *BrandHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateBrandRequest
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

	brand, err := h.brandService.Update(r.Context(), req.ID, req.Name)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	httpx.Success(w, http.StatusOK, "success update brand", brand)
}

func (h *BrandHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.HandleError(w, errors.New(errors.ErrValidation, "id is required"), h.logger)
	}

	err := h.brandService.Delete(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	httpx.Success(w, http.StatusOK, "Delete brand success", "")
}

func mapBrandToDTO(brand *sharedModel.Brand) models.BrandResponse {
	return models.BrandResponse{
		ID:        brand.ID.String(),
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
}
