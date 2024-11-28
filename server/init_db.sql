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
    password VARCHAR(255) NOT NULL
);

INSERT INTO users (username, password)
VALUES 
('testuser', '$2a$10$f.tEagZA3lZc2pUvdZrLr.N1pXgqm4fubQokrfdkdJjdWitlKDq5y')
ON CONFLICT (username) DO NOTHING;

ALTER TABLE preferences ADD CONSTRAINT unique_user_id UNIQUE (user_id);

CREATE INDEX IF NOT EXISTS idx_preferences_user_id ON preferences (user_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
