-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, password, role_id, blocked)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRoleByCode :one
SELECT * FROM roles WHERE code = $1;

-- name: GetPermissionsByRoleID :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1;

-- name: GetResourceByCode :one
SELECT * FROM resources WHERE code = $1;
