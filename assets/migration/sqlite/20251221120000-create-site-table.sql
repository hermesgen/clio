-- +migrate Up
CREATE TABLE site (
    id TEXT PRIMARY KEY,
    short_id TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    mode TEXT NOT NULL DEFAULT 'structured',
    active INTEGER NOT NULL DEFAULT 1,
    created_by TEXT NOT NULL,
    updated_by TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX idx_site_slug ON site (slug);
CREATE INDEX idx_site_active ON site (active);

-- +migrate Down
DROP TABLE site;
