-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT(user_id, role_id) DO NOTHING;

-- name: GetPrimaryRoleByUserID :one
SELECT r.name 
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1
ORDER BY ur.created_at ASC
LIMIT 1;

-- name: GetRoleByName :one
SELECT id
FROM roles
WHERE name = $1
LIMIT 1;