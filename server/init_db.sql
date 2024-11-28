CREATE TABLE IF NOT EXISTS preferences (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    sort_order TEXT NOT NULL DEFAULT 'asc',
    hidden_devices JSONB,
    latest_device_point JSONB,
    device_state JSONB,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

INSERT INTO users (username, password_hash)
VALUES 
('testuser', 'password123')
ON CONFLICT (username) DO NOTHING;

ALTER TABLE preferences ADD CONSTRAINT unique_user_id UNIQUE (user_id);

CREATE INDEX IF NOT EXISTS idx_preferences_user_id ON preferences (user_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);