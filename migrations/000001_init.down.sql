-- +migrate Down

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS roles;

DROP EXTENSION IF EXISTS "pgcrypto";
