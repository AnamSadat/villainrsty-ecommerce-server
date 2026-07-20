-- name: CreateBrand :exec
INSERT INTO
  brands (id, name, slug, created_at, updated_at)
VALUES
  ($1, $2, $3, $4, $5);


-- name: GetAllBrands :many
SELECT
  id,
  name,
  slug,
  created_at,
  updated_at
FROM
  brands
ORDER BY
  created_at DESC;


-- name: GetBrandByID :one
SELECT
  id,
  name,
  slug,
  created_at,
  updated_at
FROM
  brands
WHERE
  id = $1
LIMIT
  1;


-- name: GetBrandBySlug :one
SELECT
  id,
  name,
  slug,
  created_at,
  updated_at
FROM
  brands
WHERE
  slug = $1
LIMIT
  1;


-- name: UpdateBrandByID :exec
UPDATE
  brands
SET
  name = $2,
  slug = $3,
  updated_at = NOW()
WHERE
  id = $1;


-- name: DeleteBrandByID :exec
DELETE FROM
  brands
WHERE
  id = $1;


-- name: BrandExists :one
SELECT
  EXISTS(
    SELECT
      1
    FROM
      brands
    WHERE
      slug = $1
  );