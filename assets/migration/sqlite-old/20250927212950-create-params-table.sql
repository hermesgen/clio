-- +migrate Up
CREATE TABLE param (
    id TEXT PRIMARY KEY,
    short_id TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    value TEXT NOT NULL,
    ref_key TEXT UNIQUE,
    system INTEGER NOT NULL DEFAULT 0,
    created_by TEXT NOT NULL,
    updated_by TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX idx_param_name ON param (name);
CREATE INDEX idx_param_ref_key ON param (ref_key);

-- +migrate Down
DROP TABLE param;
