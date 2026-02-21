package handler

import (
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/lib/httpx"
)

type RbacHandler struct{}

func NewRbacHandler() *RbacHandler {
	return &RbacHandler{}
}

func (h *RbacHandler) ListProducts(w http.ResponseWriter, _ *http.Request) {
	httpx.Success(w, http.StatusOK, "products listed", map[string]any{"resource": "products", "action": "read"})
}

func (h *RbacHandler) CreateProduct(w http.ResponseWriter, _ *http.Request) {
	httpx.Success(w, http.StatusOK, "product created", map[string]any{"resource": "products", "action": "create"})
}

func (h *RbacHandler) UpdateProduct(w http.ResponseWriter, _ *http.Request) {
	httpx.Success(w, http.StatusOK, "product updated", map[string]any{"resource": "products", "action": "update"})
}

func (h *RbacHandler) DeleteProduct(w http.ResponseWriter, _ *http.Request) {
	httpx.Success(w, http.StatusOK, "product deleted", map[string]any{"resource": "products", "action": "delete"})
}

func (h *RbacHandler) ListOrders(w http.ResponseWriter, _ *http.Request) {
	httpx.Success(w, http.StatusOK, "orders listed", map[string]any{"resource": "orders", "action": "read"})
}

func (h *RbacHandler) CreateOrder(w http.ResponseWriter, _ *http.Request) {
	httpx.Success(w, http.StatusOK, "order created", map[string]any{"resource": "orders", "action": "create"})
}

func (h *RbacHandler) ManageUsers(w http.ResponseWriter, _ *http.Request) {
	httpx.Success(w, http.StatusOK, "users managed", map[string]any{"resource": "admin/users", "action": "manage"})
}
