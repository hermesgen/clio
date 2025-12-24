-- +migrate Up
CREATE TABLE meta (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL UNIQUE,
    summary TEXT NOT NULL DEFAULT '',
    excerpt TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    keywords TEXT NOT NULL DEFAULT '',
    robots TEXT NOT NULL DEFAULT '',
    canonical_url TEXT NOT NULL DEFAULT '',
    sitemap TEXT NOT NULL DEFAULT '',
    table_of_contents INTEGER NOT NULL DEFAULT 0,
    share INTEGER NOT NULL DEFAULT 0,
    comments INTEGER NOT NULL DEFAULT 0,
    created_by TEXT NOT NULL DEFAULT '',
    updated_by TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content (id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE meta;
