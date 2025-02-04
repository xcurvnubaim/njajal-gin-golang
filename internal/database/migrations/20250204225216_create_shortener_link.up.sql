CREATE TABLE shortener_links (
    id UUID PRIMARY KEY,
    original_url TEXT NOT NULL,
    shortener_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);