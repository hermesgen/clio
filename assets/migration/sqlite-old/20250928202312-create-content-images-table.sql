-- +migrate Up

CREATE TABLE content_images (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL,
    image_id TEXT NOT NULL,
    purpose TEXT NOT NULL CHECK (purpose IN ('header', 'content', 'thumbnail')),
    position INTEGER DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),

    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    UNIQUE(content_id, image_id, purpose)
);

CREATE INDEX idx_content_images_content_id ON content_images (content_id);
CREATE INDEX idx_content_images_image_id ON content_images (image_id);
CREATE INDEX idx_content_images_purpose ON content_images (purpose);
CREATE INDEX idx_content_images_active ON content_images (is_active);

-- +migrate Down

DROP TABLE content_images;