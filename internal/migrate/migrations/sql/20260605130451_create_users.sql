-- +goose Up
SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_role AS ENUM(
    'user',
    'admin'
);

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role user_role DEFAULT 'user',
    createdAt TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
SELECT 'down SQL query';

DROP TYPE IF EXISTS user_role

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS citext;
