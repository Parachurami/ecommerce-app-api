-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS products(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userId UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    budget FLOAT NOT NULL,
    skills TEXT[],
    duration INTERVAL NOT NULL,
    expiration INTERVAL NOT NULL,
    image_url TEXT,
    deliverables TEXT[],
    createdAt TIMESTAMPTZ DEFAULT NOW(),
    updatedAt TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_product FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS products;
