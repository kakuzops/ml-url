CREATE TABLE IF NOT EXISTS shorten_url (
    id VARCHAR(36) PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_shorten_url_expires_at ON shorten_url(expires_at);
CREATE INDEX idx_shorten_url_short_url ON shorten_url(short_url); 