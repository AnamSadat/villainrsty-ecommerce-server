package routes

import (
	"villainrsty-ecommerce-server/internal/adapters/http/brands/handler"

	"github.com/go-chi/chi/v5"
)

func BrandsRoute(r chi.Router, handler *handler.BrandHandler) {
	r.Route("/brands", func(r chi.Router) {
		r.Get("/by-slug/{slug}", handler.GetDetailBrandBySlug)
		r.Get("/{id}", handler.GetDetailBrandByID)
		r.Post("/", handler.AddBrand)
		r.Delete("/{id}", handler.DeleteBrand)
		r.Put("/{id}", handler.UpdateBrand)
	})
}
