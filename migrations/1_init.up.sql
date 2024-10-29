-- migrations/1_init.up.sql
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_url VARCHAR(255) NOT NULL,
    short_code VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
