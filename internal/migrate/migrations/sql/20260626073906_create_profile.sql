-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS profile(
    userId UUID PRIMARY KEY NOT NULL,
    firstName VARCHAR(255) NOT NULL,
    lastName VARCHAR(255) NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    profileImage TEXT,
    bio TEXT,
    createdAt TIMESTAMPTZ DEFAULT NOW(),
    updatedAt TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_users FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS profile;
