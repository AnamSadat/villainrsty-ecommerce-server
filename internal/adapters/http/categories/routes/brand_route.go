package routes

import (
	"villainrsty-ecommerce-server/internal/adapters/http/categories/handler"

	"github.com/go-chi/chi/v5"
)

func CategoriesRoute(r chi.Router, handler *handler.CategoryHandler) {
	r.Route("/category", func(r chi.Router) {
		r.Get("/by-slug/{slug}", handler.GetDetailCategoryBySlug)
		r.Get("/{id}", handler.GetDetailCategoryByID)
		r.Post("/", handler.AddCategory)
		r.Delete("/{id}", handler.DeleteCategory)
		r.Put("/{id}", handler.UpdateCategory)
	})
}
