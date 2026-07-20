-- name: CreateProduct :exec
INSERT INTO products (id, brand_id, category_id, name, slug, description, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetAllProducts :many
SELECT id, brand_id, category_id, name, slug, description, is_active, created_at, updated_at
FROM products
ORDER BY created_at DESC;

-- name: GetProductByID :one
SELECT id, brand_id, category_id, name, slug, description, is_active, created_at, updated_at
FROM products
WHERE id = $1
LIMIT 1;

-- name: GetProductBySlug :one
SELECT id, brand_id, category_id, name, slug, description, is_active, created_at, updated_at
FROM products
WHERE slug = $1
LIMIT 1;

-- name: ProductExists :one
SELECT EXISTS(SELECT 1 FROM products WHERE slug = $1);

-- name: UpdateProductByID :exec
UPDATE products
SET 
    brand_id = $2,
    category_id = $3,
    name = $4,
    slug = $5,
    description = $6,
    is_active = $7,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteProductByID :exec
DELETE FROM products
WHERE id = $1;