-- +migrate Up
-- Note: Site directories need to be created separately by the workspace setup
-- This migration only inserts the site records

-- Insert structured site
INSERT OR IGNORE INTO site (id, short_id, name, slug, mode, active, created_by, updated_by, created_at, updated_at)
VALUES (
	'a0000000-0000-0000-0000-000000000001',
	'strc',
	'Structured Site',
	'structured',
	'structured',
	1,
	'a0000000-0000-0000-0000-000000000001',
	'a0000000-0000-0000-0000-000000000001',
	datetime('now'),
	datetime('now')
);

-- Insert blog site
INSERT OR IGNORE INTO site (id, short_id, name, slug, mode, active, created_by, updated_by, created_at, updated_at)
VALUES (
	'b0000000-0000-0000-0000-000000000001',
	'blog',
	'Blog',
	'blog',
	'blog',
	1,
	'b0000000-0000-0000-0000-000000000001',
	'b0000000-0000-0000-0000-000000000001',
	datetime('now'),
	datetime('now')
);

-- +migrate Down
DELETE FROM site WHERE slug IN ('structured', 'blog');
