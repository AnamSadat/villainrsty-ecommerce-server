package routes

import (
	"villainrsty-ecommerce-server/internal/adapters/http/rbac/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterRbacRoute(r chi.Router, h *handler.RbacHandler) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.ListProducts)
		r.Post("/", h.CreateProduct)
		r.Patch("/{id}", h.UpdateProduct)
		r.Delete("/{id}", h.DeleteProduct)
	})

	r.Route("/orders", func(r chi.Router) {
		r.Get("/", h.ListOrders)
		r.Post("/", h.CreateOrder)
	})

	r.Route("/admin/users", func(r chi.Router) {
		r.Get("/", h.ManageUsers)
		r.Post("/", h.ManageUsers)
		r.Delete("/{id}", h.ManageUsers)
	})
}
