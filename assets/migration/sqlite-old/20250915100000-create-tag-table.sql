-- +migrate Up
CREATE TABLE tag (
    id TEXT PRIMARY KEY,
    short_id TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE tag;
