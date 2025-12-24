-- +migrate Up
CREATE TABLE IF NOT EXISTS user_site (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	site_id TEXT NOT NULL,
	role TEXT DEFAULT 'editor',
	created_at TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
	FOREIGN KEY (site_id) REFERENCES site(id) ON DELETE CASCADE,
	UNIQUE(user_id, site_id)
);

CREATE INDEX IF NOT EXISTS idx_user_site_user_id ON user_site(user_id);
CREATE INDEX IF NOT EXISTS idx_user_site_site_id ON user_site(site_id);

-- +migrate Down
DROP TABLE IF EXISTS user_site;
