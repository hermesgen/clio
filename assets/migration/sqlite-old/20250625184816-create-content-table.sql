-- +migrate Up
CREATE TABLE content (
    id TEXT PRIMARY KEY,
    short_id TEXT NOT NULL DEFAULT '',
    user_id TEXT NOT NULL,
    section_id TEXT NOT NULL,
    kind TEXT NOT NULL DEFAULT 'article',
    image TEXT NOT NULL DEFAULT '',
    heading TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    draft INTEGER NOT NULL DEFAULT 1,
    featured INTEGER NOT NULL DEFAULT 0,
    series TEXT NOT NULL DEFAULT '',
    series_order INTEGER NOT NULL DEFAULT 0,
    published_at TIMESTAMP,
    created_by TEXT NOT NULL DEFAULT '',
    updated_by TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE content;