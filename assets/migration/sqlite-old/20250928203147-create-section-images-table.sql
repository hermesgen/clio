-- +migrate Up

CREATE TABLE section_images (
    id TEXT PRIMARY KEY,
    section_id TEXT NOT NULL,
    image_id TEXT NOT NULL,
    purpose TEXT NOT NULL CHECK (purpose IN ('header', 'blog_header')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),

    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    UNIQUE(section_id, image_id, purpose)
);

CREATE INDEX idx_section_images_section_id ON section_images (section_id);
CREATE INDEX idx_section_images_image_id ON section_images (image_id);
CREATE INDEX idx_section_images_purpose ON section_images (purpose);
CREATE INDEX idx_section_images_active ON section_images (is_active);

-- +migrate Down
DROP TABLE section_images;