-- +migrate Up
INSERT INTO site (id, short_id, name, slug, mode, active, created_by, updated_by, created_at, updated_at)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'dflt',
    'Default Site',
    'default',
    'structured',
    1,
    'a0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    datetime('now'),
    datetime('now')
);

-- +migrate Down
DELETE FROM site WHERE slug = 'default';
