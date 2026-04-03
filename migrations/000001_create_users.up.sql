-- Users table based on the User proto definition
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at BIGINT DEFAULT EXTRACT(epoch FROM NOW())::BIGINT,
    updated_at BIGINT DEFAULT EXTRACT(epoch FROM NOW())::BIGINT,
    deleted_at BIGINT DEFAULT NULL
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
