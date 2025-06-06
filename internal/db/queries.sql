-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, password, role_id, blocked)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
    password = $3,
    role_id = $4,
    blocked = $5
WHERE id = $1
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
-- name: GetRoleByCode :one
SELECT * FROM roles WHERE code = $1;

-- name: GetPermissionsByRoleID :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1;

-- name: GetResourceByCode :one
SELECT * FROM resources WHERE code = $1;

-- name: CreateRole :one
INSERT INTO roles (id, code)
VALUES ($1, $2)
RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1;