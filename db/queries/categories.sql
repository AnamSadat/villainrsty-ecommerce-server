-- name: CreateCategory :exec
INSERT INTO categories (id, name, slug, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetAllCategories :many
SELECT id, name, slug, created_at, updated_at
FROM categories
ORDER BY created_at DESC;

-- name: GetCategoryByID :one
SELECT id, name, slug, created_at, updated_at
FROM categories
WHERE id = $1
LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT id, name, slug, created_at, updated_at
FROM categories
WHERE slug = $1
LIMIT 1;

-- name: UpdateCategoryByID :exec
UPDATE categories
SET name = $2, slug = $3, updated_at = NOW()
WHERE id = $1;

-- name: DeleteCategoryByID :exec
DELETE FROM categories
WHERE id = $1;

-- name: CategoryExists :one
SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1);